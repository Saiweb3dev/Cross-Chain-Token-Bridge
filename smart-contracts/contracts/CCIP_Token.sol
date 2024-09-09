// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title CCIPToken
 * @dev Implementation of the CCIPToken
 * This contract extends ERC20, Ownable, and Pausable functionalities.
 */
contract CCIP_Token is ERC20, Ownable, Pausable {
    // Custom error declarations
    error InvalidAmount(uint256 amount);
    error InsufficientBalance(address account, uint256 balance, uint256 required);
    error UnauthorizedAccess(address caller);
    error ZeroAddress();
    error InsufficientSupply(uint256 available, uint256 requested);

    // State variables
    uint256 private _tokenSupply;
    uint256 private _availableSupply;
    address private vaultAddress;

    // Named Mappings
    mapping(address sender => mapping(address recipient => uint256 amount)) public transferAmounts;
    mapping(address sender => address[] recipients) public transferAddresses;
    mapping(address user => uint256 amount) public lockedAmounts;
    
    // Events
    event Mint(address indexed to, uint256 amount);
    event Burn(address indexed from, uint256 amount);
    event CurrentOwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event TokenSupplyIncreased(uint256 amount);

    /**
     * @dev Constructor that sets up the CCIPToken.
     * @param name The name of the token.
     * @param symbol The symbol of the token.
     * @param __tokenSupply The initial token supply.
     */
    constructor(string memory name, string memory symbol, uint256 __tokenSupply)
        ERC20(name, symbol)
    {
        _tokenSupply = __tokenSupply;
        _availableSupply = __tokenSupply;
    }

    // Modifiers
    modifier validAmount(uint256 amount) {
        if (amount == 0) {
            revert InvalidAmount(amount);
        }
        _;
    }

    modifier onlyVault() {
        require(msg.sender == vaultAddress, "Caller is not the vault");
        _;
    }

    // External and public functions (non-view)

    /**
     * @dev Mints new tokens.
     * @param to The address that will receive the minted tokens.
     * @param amount The amount of tokens to mint.
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
        emit Mint(to, amount);
    }

    /**
     * @dev Burns tokens.
     * @param from The address from which to burn tokens.
     * @param amount The amount of tokens to burn.
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
        _availableSupply += amount;
        emit Burn(from, amount);
    }

    /**
     * @dev Increases the token supply.
     * @param amount The amount to increase the supply by.
     */
    function increaseTokenSupply(uint256 amount) public onlyOwner {
        _tokenSupply += amount;
        _availableSupply += amount;
        emit TokenSupplyIncreased(amount);
    }

    /**
     * @dev Transfers tokens to a specified address.
     * @param recipient The address to transfer to.
     * @param amount The amount to be transferred.
     * @return A boolean that indicates if the operation was successful.
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
     * @dev Transfers tokens to the CCIP vault.
     * @param _amount The amount to transfer.
     */
    function transferToCCIPVault(uint256 _amount) public {
        _transfer(msg.sender, vaultAddress, _amount);
        lockedAmounts[msg.sender] += _amount;
    }

    /**
     * @dev Burns the locked amount for a user.
     * @param _user The address of the user.
     * @param _amount The amount to burn.
     */
    function burnLockAmount(address _user, uint256 _amount) public onlyVault {
        require(lockedAmounts[_user] >= _amount, "Insufficient locked amount");
        lockedAmounts[_user] -= _amount;
    }

    /**
     * @dev Sets the vault address.
     * @param _newVaultAddress The new vault address.
     */
    function setVaultAddress(address _newVaultAddress) public onlyOwner {
        require(_newVaultAddress != address(0), "Invalid address");
        vaultAddress = _newVaultAddress;
    }

    /**
     * @dev Transfers tokens from one address to another.
     * @param sender The address to transfer from.
     * @param recipient The address to transfer to.
     * @param amount The amount to be transferred.
     * @return A boolean that indicates if the operation was successful.
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
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * @param newOwner The address of the new owner.
     */
    function transferOwnership(address newOwner) public virtual override onlyOwner {
        if (newOwner == address(0)) {
            revert ZeroAddress();
        }
        _transferOwnership(newOwner);
        emit CurrentOwnershipTransferred(owner(), newOwner);
    }

    // Public view functions

    /**
     * @dev Returns the initial supply of tokens.
     * @return The initial token supply.
     */
    function initialSupply() public view returns (uint256) {
        return _tokenSupply;
    }

    /**
     * @dev Returns the available supply of tokens.
     * @return The available token supply.
     */
    function availableSupply() public view returns (uint256) {
        return _availableSupply;
    }

    /**
     * @dev Returns the total supply of tokens.
     * @return The total token supply.
     */
    function totalSupply() public view override returns (uint256) {
        return _tokenSupply - _availableSupply;
    }

    /**
     * @dev Returns the balance of the given account.
     * @param account The address to query the balance of.
     * @return The balance of the given account.
     */
    function balanceOf(address account) public view override returns (uint256) {
        return super.balanceOf(account);
    }

    /**
     * @dev Returns the locked amount for a user.
     * @param _user The address of the user.
     * @return The locked amount.
     */
    function getLockAmount(address _user) public view returns(uint256) {
        return lockedAmounts[_user];
    }
}