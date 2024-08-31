// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

// Custom error messages for better error handling
error InvalidAmount(uint256 amount);
error InsufficientBalance(address account, uint256 balance, uint256 required);
error UnauthorizedAccess(address caller);
error ZeroAddress();

/**
 * @title CCIPToken
 * @dev An ERC20 token with additional functionality like pausability and ownership control
 */
contract CCIPToken is ERC20, Ownable, Pausable {
    // Events emitted during token operations
    event Mint(address indexed to, uint256 amount);
    event Burn(address indexed from, uint256 amount);

    /**
     * @dev Constructor initializes the token with given name and symbol
     * @param name Token name
     * @param symbol Token symbol
     */
    constructor(string memory name, string memory symbol)
        ERC20(name, symbol)
    {}

    /**
     * @dev Modifier to validate the amount before minting
     * @param amount Amount to validate
     */
    modifier validAmount(uint256 amount) {
        if (amount == 0 || amount > type(uint256).max - totalSupply()) {
            revert InvalidAmount(amount);
        }
        _;
    }

    /**
     * @dev Mint tokens to a specified address
     * @param to Recipient address
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) public onlyOwner whenNotPaused validAmount(amount) {
        if (to == address(0)) {
            revert ZeroAddress();
        }
        _mint(to, amount);
        emit Mint(to, amount);
    }

    /**
     * @dev Burn tokens from a specified address
     * @param from Sender address
     * @param amount Amount to burn
     */
    function burn(address from, uint256 amount) public onlyOwner whenNotPaused validAmount(amount) {
        if (from == address(0)) {
            revert ZeroAddress();
        }
        if (balanceOf(from) < amount) {
            revert InsufficientBalance(from, balanceOf(from), amount);
        }
        _burn(from, amount);
        emit Burn(from, amount);
    }

    /**
     * @dev Pause token transfers
     */
    function pause() external onlyOwner {
        _pause();
    }

    /**
     * @dev Unpause token transfers
     */
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Get the decimal places of the token
     * @return Decimals of the token
     */
    function decimals() public pure override returns (uint8) {
        return 18;
    }

    /**
     * @dev Get the total supply of the token
     * @return Total supply of the token
     */
    function totalSupply() public view override returns (uint256) {
        return super.totalSupply();
    }

    /**
     * @dev Get the balance of a specific account
     * @param account Address to check balance for
     * @return Balance of the account
     */
    function balanceOf(address account) public view override returns (uint256) {
        return super.balanceOf(account);
    }

    /**
     * @dev Transfer tokens between accounts
     * @param recipient Recipient address
     * @param amount Amount to transfer
     * @return Success flag
     */
    function transfer(address recipient, uint256 amount) public override whenNotPaused returns (bool) {
        _transfer(msg.sender, recipient, amount);
        return true;
    }

    /**
     * @dev Transfer tokens from one account to another
     * @param sender Sender address
     * @param recipient Recipient address
     * @param amount Amount to transfer
     * @return Success flag
     */
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

    /**
     * @dev Get the owner address of the contract
     * @return Owner address
     */
    function OwnerAddress() public view returns (address) {
        return owner();
    }
}
