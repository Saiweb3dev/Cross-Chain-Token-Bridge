// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Vault is ReentrancyGuard, Ownable {

    mapping(address => mapping(address => uint256)) public transferAmounts;
    mapping(address => address[]) public transferAddresses;

    event countNumberFromSourceEvent(address indexed from,address indexed to, uint256 amount);
    event countNumberFromDestinationEvent(address indexed from,address indexed to, uint256 amount);
    event TokensReleased(address indexed user, uint256 amount);

    error InvalidAmount(uint256 amount);
    error InsufficientBalance(address account, uint256 balance, uint256 required);

    constructor() Ownable(msg.sender) {
    }


    function countNumberFromSource(
        address _from,
        address _to,
        uint256 _amount
    ) external {
        transferAmounts[_from][_to] += _amount;
        if (transferAmounts[_from][_to] > 0 && !containsAddress(_from, _to)) {
            transferAddresses[_from].push(_to);
        }
        emit countNumberFromSourceEvent(_from,_to,_amount);
    }
    function countNumberFromDestination(
        address _from,
        address _to,
        uint256 _amount
    ) external {
        transferAmounts[_from][_to] += _amount;
        if (transferAmounts[_from][_to] > 0 && !containsAddress(_from, _to)) {
            transferAddresses[_from].push(_to);
        }
        emit countNumberFromDestinationEvent(_from,_to,_amount);
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