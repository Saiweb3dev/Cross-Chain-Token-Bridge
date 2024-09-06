// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

interface Token {
    function mint(address to, uint256 amount) external;

    function burn(address from, uint256 amount) external;

    function getLockAmount(address _user) external returns (uint256);

    function balanceOf(address account) external view returns (uint256);
}

error AmountDoestMatch();
error AmountDoesNotMatch(uint256 expected, uint256 actual);
error InvalidAddress();

contract Vault is ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    mapping(address => mapping(address => uint256)) public transferAmounts;
    mapping(address => address[]) public transferAddresses;

    event CountNumberFromSourceEvent(
        address indexed from,
        address indexed to,
        uint256 amount
    );
    event CountNumberFromDestinationEvent(
        address indexed from,
        address indexed to,
        uint256 amount
    );
    event TokensReleased(address indexed user, uint256 amount);
    event tokenLocked(address indexed _from, uint256 amount);
    event TokensLocked(address indexed user, uint256 amount);

    error InvalidAmount(uint256 amount);
    error InsufficientBalance(uint256 account, uint256 balance);

    Token private token;
    

    constructor(address _token) Ownable(msg.sender) {
        require(_token != address(0), "Token address cannot be zero");
        token = Token(_token);
    }

    function lockTokenInVault(address _from, uint256 _amount) external {
        // Input validation
        if (_from == address(0)) revert InvalidAddress();
        if (_amount == 0)
            revert InvalidAmount(_amount);

        // Get locked amount from token contract
        uint256 lockedAmount = token.getLockAmount(_from);

        // Compare amounts (using == instead of >=)
        // if (lockedAmount != _amount) {
        //     revert AmountDoesNotMatch(_amount, lockedAmount);
        // }

        // Check vault balance
        uint256 vaultBalance = token.balanceOf(address(this));
        // if (vaultBalance < _amount) {
        //     revert InsufficientBalance(vaultBalance, _amount);
        // }

        // Burn tokens
        token.burn(address(this), _amount);

        emit TokensLocked(_from, _amount);
    }

    function releaseTokenInVault(address _to, uint256 _amount) external {
        require(_to != address(0), "To address cannot be zero");
        require(_amount > 0, "Amount must be greater than zero");
        token.mint(_to, _amount);
    }

    function countNumberFromSource(
        address _from,
        address _to,
        uint256 _amount
    ) external {
        require(
            _from != address(0) && _to != address(0),
            "Addresses cannot be zero"
        );
        require(_amount > 0, "Amount must be greater than zero");

        transferAmounts[_from][_to] += _amount;
        if (transferAmounts[_from][_to] > 0 && !containsAddress(_from, _to)) {
            transferAddresses[_from].push(_to);
        }

        emit CountNumberFromSourceEvent(_from, _to, _amount);
    }

    function countNumberFromDestination(
        address _from,
        address _to,
        uint256 _amount
    ) external {
        require(
            _from != address(0) && _to != address(0),
            "Addresses cannot be zero"
        );
        require(_amount > 0, "Amount must be greater than zero");

        transferAmounts[_from][_to] += _amount;
        if (transferAmounts[_from][_to] > 0 && !containsAddress(_from, _to)) {
            transferAddresses[_from].push(_to);
        }
        emit CountNumberFromDestinationEvent(_from, _to, _amount);
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
}
