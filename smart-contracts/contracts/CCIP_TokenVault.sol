// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// Custom interface extending IERC20 with burn and mint functions
interface ICCIPToken is IERC20 {
    function burn(address from, uint256 amount) external;
    function mint(address to, uint256 amount) external;
}

/**
 * @title CCIPvault
 * @dev A vault contract for managing token locks and releases
 */
contract CCIPvault is ReentrancyGuard, Ownable {
    using SafeERC20 for ICCIPToken;

    // The token being managed by this vault
    ICCIPToken public immutable token;
    
    // Flag to track whether the vault is paused
    bool public paused;

    // Mapping of user addresses to their locked balances
    mapping(address => uint256) public lockedBalances;

    // Events emitted during token operations
    event TokensLocked(address indexed user, uint256 amount);
    event TokensReleased(address indexed user, uint256 amount);

    // Custom error messages for better error handling
    error InvalidAmount(uint256 amount);
    error InsufficientBalance(address account, uint256 balance, uint256 required);

    /**
     * @dev Constructor initializes the vault with the given token address
     * @param _tokenAddress Address of the token to manage
     */
    constructor(address _tokenAddress){
        token = ICCIPToken(_tokenAddress);
    }

    /**
     * @dev Modifier to ensure the contract is not paused
     */
    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    /**
     * @dev Locks tokens from a user and burns them
     * @param _from User address to lock tokens from
     * @param _amount Amount of tokens to lock
     * @return bool Whether the operation was successful
     */
    function lockTokens(address _from, uint256 _amount) external whenNotPaused nonReentrant returns (bool) {
        if (_amount == 0) {
            revert InvalidAmount(_amount);
        }
        if (token.balanceOf(_from) < _amount) {
            revert InsufficientBalance(_from, token.balanceOf(_from), _amount);
        }

        token.safeTransferFrom(_from, address(this), _amount);
        lockedBalances[_from] += _amount;
        emit TokensLocked(_from, _amount);

        // Burn tokens
        token.burn(address(this), _amount);

        return true;
    }

    /**
     * @dev Releases locked tokens back to a user and mints new ones
     * @param _to User address to receive released tokens
     * @param _amount Amount of tokens to release
     * @return bool Whether the operation was successful
     */
    function releaseTokens(address _to, uint256 _amount) external whenNotPaused nonReentrant returns (bool) {
        if (_amount == 0) {
            revert InvalidAmount(_amount);
        }

        // Mint new tokens
        token.mint(_to, _amount);
        emit TokensReleased(_to, _amount);

        return true;
    }

    /**
     * @dev Pause all operations on the vault
     */
    function pause() external onlyOwner {
        paused = true;
    }

    /**
     * @dev Unpause all operations on the vault
     */
    function unpause() external onlyOwner {
        paused = false;
    }

    /**
     * @dev Get the total balance of tokens held by the vault
     * @return uint256 Total token balance
     */
    function getTokenBalance() external view returns (uint256) {
        return token.balanceOf(address(this));
    }

    /**
     * @dev Get the locked balance of a specific user
     * @param _user User address to check
     * @return uint256 Locked balance of the user
     */
    function getUserLockedBalance(address _user) external view returns (uint256) {
        return lockedBalances[_user];
    }
}
