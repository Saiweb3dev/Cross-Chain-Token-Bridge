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

    // Mapping to track transfer amounts between addresses
    mapping(address => mapping(address => uint256)) public transferAmounts;
    // Mapping to track transfer addresses for each address
    mapping(address => address[]) public transferAddresses;

    // Events
    event CountNumberFromSourceEvent(address indexed from, address indexed to, uint256 amount);
    event CountNumberFromDestinationEvent(address indexed from, address indexed to, uint256 amount);
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
    constructor(address _token) Ownable(msg.sender) {
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

    /**
     * @dev Checks if an address is contained in the transfer addresses of another address
     * @param _address The main address
     * @param _otherAddress The address to check for
     * @return bool Returns true if _otherAddress is in the transfer addresses of _address
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
     * @dev Gets the total transfer count from an address
     * @param _address The address to check
     * @return uint256 The total transfer count
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
}