// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title TokenVault
 * @dev A simple token vault for managing ERC-20 token deposits and withdrawals
 */
contract TokenVault is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // The ERC-20 token managed by this vault
    IERC20 public immutable token;

    // Indicates if the contract is paused
    bool public paused;

    // Mapping of user addresses to their token balances
    mapping(address => uint256) public userBalances;

    // Events
    event Deposit(address indexed user, uint256 amount);
    event Withdrawal(address indexed user, uint256 amount);
    event EmergencyWithdrawal(address indexed user, uint256 amount);

    /**
     * @dev Sets the token address and transfers ownership to the deployer
     * @param _tokenAddress The address of the ERC-20 token to be managed
     */
    constructor(address _tokenAddress){
        token = IERC20(_tokenAddress);
    }

    /**
     * @dev Ensures the contract is not paused
     */
    modifier whenNotPaused() {
        require(!paused, "Contract is paused");
        _;
    }

    /**
     * @dev Allows users to deposit tokens
     * @param _amount The amount of tokens to deposit
     */
    function deposit(uint256 _amount) external nonReentrant whenNotPaused {
        require(_amount > 0, "Deposit amount must be greater than zero");
        
        token.safeTransferFrom(msg.sender, address(this), _amount);
        userBalances[msg.sender] += _amount;
        
        emit Deposit(msg.sender, _amount);
    }

    /**
     * @dev Allows users to withdraw tokens
     * @param _amount The amount of tokens to withdraw
     */
    function withdraw(uint256 _amount) external nonReentrant whenNotPaused {
        require(_amount > 0, "Withdrawal amount must be greater than zero");
        require(userBalances[msg.sender] >= _amount, "Insufficient balance");

        userBalances[msg.sender] -= _amount;
        token.safeTransfer(msg.sender, _amount);

        emit Withdrawal(msg.sender, _amount);
    }

    /**
     * @dev Allows emergency withdrawal of all user funds
     */
    function emergencyWithdraw() external nonReentrant {
        uint256 amount = userBalances[msg.sender];
        require(amount > 0, "No balance to withdraw");
        
        userBalances[msg.sender] = 0;
        token.safeTransfer(msg.sender, amount);
        
        emit EmergencyWithdrawal(msg.sender, amount);
    }

    /**
     * @dev Allows batch deposits for multiple users
     * @param users Array of user addresses
     * @param amounts Array of deposit amounts
     */
    function batchDeposit(address[] calldata users, uint256[] calldata amounts) external nonReentrant whenNotPaused {
        require(users.length == amounts.length, "Arrays length mismatch");
        
        for (uint256 i = 0; i < users.length; i++) {
            uint256 amount = amounts[i];
            token.safeTransferFrom(msg.sender, address(this), amount);
            userBalances[users[i]] += amount;
            
            emit Deposit(users[i], amount);
        }
    }

    /**
     * @dev Pauses the contract
     */
    function pause() external onlyOwner {
        paused = true;
    }

    /**
     * @dev Unpauses the contract
     */
    function unpause() external onlyOwner {
        paused = false;
    }

    /**
     * @dev Returns the total balance of tokens in the vault
     * @return Total balance of tokens
     */
    function getTotalBalance() external view returns (uint256) {
        return token.balanceOf(address(this));
    }

    /**
     * @dev Returns the balance of a specific user
     * @param _user Address of the user
     * @return Balance of the user
     */
    function getUserBalance(address _user) external view returns (uint256) {
        return userBalances[_user];
    }
}