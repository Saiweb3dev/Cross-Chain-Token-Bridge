// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

error InvalidAmount(uint256 amount);
error InsufficientBalance(address account, uint256 balance, uint256 required);
error UnauthorizedAccess(address caller);
error ZeroAddress();
error InsufficientSupply(uint256 available, uint256 requested);

contract CCIPToken is ERC20, Ownable, Pausable {
     uint256 private _initialSupply;
    uint256 private _availableSupply;
    uint256 private _tokenInSupply;
    address private vaultAddress;

    mapping(address => mapping(address => uint256)) public transferAmounts;
    mapping(address => address[]) public transferAddresses;
    mapping(address => uint256) public lockedAmounts;
    
    event Mint(address indexed to, uint256 amount);
    event Burn(address indexed from, uint256 amount);
    event CurrentOwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    constructor(string memory name, string memory symbol, uint256 __initialSupply)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {
        _initialSupply = __initialSupply;
        _availableSupply = __initialSupply;
        _mint(msg.sender, __initialSupply);
    }

    modifier validAmount(uint256 amount) {
        if (amount == 0) {
            revert InvalidAmount(amount);
        }
        _;
    }

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

    function initialSupply() public view returns (uint256) {
        return _initialSupply;
    }

    function availableSupply() public view returns (uint256) {
        return _availableSupply;
    }

    function tokenInSupply() public view returns (uint256) {
        return _tokenInSupply;
    }

    function balanceOf(address account) public view override returns (uint256) {
        return super.balanceOf(account);
    }

    function transfer(address recipient, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    function transferToCCIPVault(uint256 _amount) public {
        transferFrom(msg.sender,vaultAddress,_amount);
        lockedAmounts[msg.sender] += _amount;
        
    }

    function getLockAmount(address _user) public view returns(uint256){
        return lockedAmounts[_user];
    }

    function updateVaultAddress(address _vaultAddress) public {
        vaultAddress = _vaultAddress;
    }

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public override whenNotPaused returns (bool) {
        _transfer(sender, recipient, amount);

        uint256 currentAllowance = allowance(sender, msg.sender);
        require(currentAllowance >= amount, "ERC20: insufficient allowance");
        unchecked {
            _approve(sender, msg.sender, currentAllowance - amount);
        }

        return true;
    }

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

    function OwnerAddress() public view returns (address) {
        return owner();
    }
     function transferOwnership(address newOwner) public virtual override onlyOwner {
        if (newOwner == address(0)) {
            revert ZeroAddress();
        }
        _transferOwnership(newOwner);
        emit CurrentOwnershipTransferred(owner(), newOwner);
    }
}
