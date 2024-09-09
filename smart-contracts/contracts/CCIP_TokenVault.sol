// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title Token Interface
 * @dev Interface for interacting with the CCIPToken contract
 */
interface Token {
    function mint(address to, uint256 amount) external;
    function burn(address from, uint256 amount) external;
    function getLockAmount(address _user) external returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function burnLockAmount(address _user, uint256 _amount) external;
}

// Custom error definitions
error AmountDoestMatch();
error AmountDoesNotMatch(uint256 expected, uint256 actual);
error InvalidAddress();

/**
 * @title Vault
 * @dev A contract for managing token locking and releasing operations
 * This contract implements ReentrancyGuard and Ownable for security and access control
 */
contract Vault is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // Events
    event TokensReleased(address indexed user, uint256 amount);
    event tokenLocked(address indexed _from, uint256 amount);
    event TokensLocked(address indexed user, uint256 amount);

    // Custom errors
    error InvalidAmount(uint256 amount);
    error InsufficientBalance(uint256 account, uint256 balance);

    // Token contract instance
    Token private token;

    /**
     * @dev Constructor that sets the address of the token contract
     * @param _token Address of the Token contract
     */
    constructor(address _token) {
        require(_token != address(0), "Token address cannot be zero");
        token = Token(_token);
    }

    /**
     * @dev Locks tokens in the vault
     * @param _from Address from which tokens are being locked
     * @param _amount Amount of tokens to lock
     */
    function lockTokenInVault(address _from, uint256 _amount) external {
        // Input validation
        if (_from == address(0)) revert InvalidAddress();
        if (_amount == 0) revert InvalidAmount(_amount);

        // Check vault balance
        uint256 vaultBalance = token.balanceOf(address(this));
        if (vaultBalance < _amount) {
            revert InsufficientBalance(vaultBalance, _amount);
        }

        // Burn tokens
        token.burn(address(this), _amount);
        token.burnLockAmount(_from, _amount);

        emit TokensLocked(_from, _amount);
    }

    /**
     * @dev Releases tokens from the vault
     * @param _to Address to which tokens are being released
     * @param _amount Amount of tokens to release
     */
    function releaseTokenInVault(address _to, uint256 _amount) external {
        require(_to != address(0), "To address cannot be zero");
        require(_amount > 0, "Amount must be greater than zero");
        token.mint(_to, _amount);
    }

}