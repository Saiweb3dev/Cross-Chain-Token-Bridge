// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

// Custom error definitions
error InvalidAmount(uint256 amount);
error InsufficientBalance(address account, uint256 balance, uint256 required);
error UnauthorizedAccess(address caller);
error ZeroAddress();
error InsufficientSupply(uint256 available, uint256 requested);

/**
 * @title CCIPToken
 * @dev Implementation of the CCIPToken
 * This contract extends ERC20, Ownable, and Pausable functionalities.
 */
contract CCIPToken is ERC20, Ownable, Pausable {
    uint256 private _initialSupply;
    uint256 private _availableSupply;
    uint256 private _tokenInSupply;
    address private vaultAddress;

    // Mapping to track transfer amounts between addresses
    mapping(address => mapping(address => uint256)) public transferAmounts;
    // Mapping to track transfer addresses for each address
    mapping(address => address[]) public transferAddresses;
    // Mapping to track locked amounts for each address
    mapping(address => uint256) public lockedAmounts;
    
    // Events
    event Mint(address indexed to, uint256 amount);
    event Burn(address indexed from, uint256 amount);
    event CurrentOwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /**
     * @dev Constructor that gives msg.sender all of existing tokens.
     * @param name The name of the token
     * @param symbol The symbol of the token
     * @param __initialSupply The initial supply of tokens
     */
    constructor(string memory name, string memory symbol, uint256 __initialSupply)
        ERC20(name, symbol)
    {
        _initialSupply = __initialSupply;
        _availableSupply = __initialSupply;
        _mint(msg.sender, __initialSupply);
    }

    /**
     * @dev Modifier to validate if the amount is greater than zero
     */
    modifier validAmount(uint256 amount) {
        if (amount == 0) {
            revert InvalidAmount(amount);
        }
        _;
    }

    /**
     * @dev Modifier to restrict access to only the vault address
     */
    modifier onlyVault() {
        require(msg.sender == vaultAddress, "Caller is not the vault");
        _;
    }

    /**
     * @dev Function to mint tokens
     * @param to The address that will receive the minted tokens
     * @param amount The amount of tokens to mint
     */
    function mint(address to, uint256 amount)
        public
        validAmount(amount)
    {
        if (to == address(0)) {
            revert ZeroAddress();
        }
        if (amount > _availableSupply) {
            revert InsufficientSupply(_availableSupply, amount);
        }
        _mint(to, amount);
        _availableSupply -= amount;
        _tokenInSupply += amount;
        emit Mint(to, amount);
    }

    /**
     * @dev Function to burn tokens
     * @param from The address from which to burn tokens
     * @param amount The amount of tokens to burn
     */
    function burn(address from, uint256 amount)
        public
        validAmount(amount)
    {
        if (from == address(0)) {
            revert ZeroAddress();
        }
        if (balanceOf(from) < amount) {
            revert InsufficientBalance(from, balanceOf(from), amount);
        }
        _burn(from, amount);
        _tokenInSupply -= amount;
        emit Burn(from, amount);
    }

    /**
     * @dev Returns the initial supply of tokens
     */
    function initialSupply() public view returns (uint256) {
        return _initialSupply;
    }

    /**
     * @dev Returns the available supply of tokens
     */
    function availableSupply() public view returns (uint256) {
        return _availableSupply;
    }

    /**
     * @dev Returns the current supply of tokens in circulation
     */
    function tokenInSupply() public view returns (uint256) {
        return _tokenInSupply;
    }

    /**
     * @dev Returns the balance of the specified account
     * @param account The address to query the balance of
     */
    function balanceOf(address account) public view override returns (uint256) {
        return super.balanceOf(account);
    }

    /**
     * @dev Transfer token to a specified address
     * @param recipient The address to transfer to
     * @param amount The amount to be transferred
     */
    function transfer(address recipient, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    /**
     * @dev Transfer tokens to the CCIP vault
     * @param _amount The amount of tokens to transfer
     */
    function transferToCCIPVault(uint256 _amount) public {
        _transfer(msg.sender, vaultAddress, _amount);
        lockedAmounts[msg.sender] += _amount;
    }

    /**
     * @dev Get the locked amount for a specific user
     * @param _user The address of the user
     */
    function getLockAmount(address _user) public view returns(uint256) {
        return lockedAmounts[_user];
    }

    /**
     * @dev Burn locked amount for a specific user
     * @param _user The address of the user
     * @param _amount The amount to burn
     */
    function burnLockAmount(address _user, uint256 _amount) public onlyVault {
        require(lockedAmounts[_user] >= _amount, "Insufficient locked amount");
        lockedAmounts[_user] -= _amount;
    }

    /**
     * @dev Set the vault address
     * @param _newVaultAddress The new vault address
     */
    function setVaultAddress(address _newVaultAddress) public onlyOwner {
        require(_newVaultAddress != address(0), "Invalid address");
        vaultAddress = _newVaultAddress;
    }

    /**
     * @dev Transfer tokens from one address to another
     * @param sender The address which you want to send tokens from
     * @param recipient The address which you want to transfer to
     * @param amount The amount of tokens to be transferred
     */
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public override whenNotPaused returns (bool) {
        _transfer(sender, recipient, amount);
        return true;
    }

    /**
     * @dev Count the number of transfers between addresses
     * @param _from The sender's address
     * @param _to The recipient's address
     * @param _amount The amount transferred
     */
    function countNumber(
        address _from,
        address _to,
        uint256 _amount
    ) external {
        transferAmounts[_from][_to] += _amount;
        if (transferAmounts[_from][_to] > 0 && !containsAddress(_from, _to)) {
            transferAddresses[_from].push(_to);
        }
    }

    /**
     * @dev Check if an address is contained in the transfer addresses of another address
     * @param _address The main address
     * @param _otherAddress The address to check for
     */
    function containsAddress(address _address, address _otherAddress)
        internal
        view
        returns (bool)
    {
        address[] memory addresses = transferAddresses[_address];
        for (uint256 i = 0; i < addresses.length; i++) {
            if (addresses[i] == _otherAddress) return true;
        }
        return false;
    }

    /**
     * @dev Get the total transfer count from an address
     * @param _address The address to check
     */
    function getTransferCountFromAddress(address _address)
        public
        view
        returns (uint256)
    {
        uint256 totalCount = 0;
        address[] memory toAddresses = transferAddresses[_address];

        for (uint256 i = 0; i < toAddresses.length; i++) {
            totalCount += transferAmounts[_address][toAddresses[i]];
        }

        return totalCount;
    }

    /**
     * @dev Returns the address of the current owner
     */
    function OwnerAddress() public view returns (address) {
        return owner();
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Can only be called by the current owner.
     * @param newOwner The address of the new owner
     */
    function transferOwnership(address newOwner) public virtual override onlyOwner {
        if (newOwner == address(0)) {
            revert ZeroAddress();
        }
        _transferOwnership(newOwner);
        emit CurrentOwnershipTransferred(owner(), newOwner);
    }
}