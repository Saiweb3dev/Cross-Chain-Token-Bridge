// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "@chainlink/contracts-ccip/src/v0.8/ccip/interfaces/IRouterClient.sol";
import {OwnerIsCreator} from "@chainlink/contracts-ccip/src/v0.8/shared/access/OwnerIsCreator.sol";
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol";
import {IERC20} from "@chainlink/contracts-ccip/src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@chainlink/contracts-ccip/src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";


/// @title Vault Interface
/// @notice Interface for interacting with a token vault
interface CCIP_TokenVault {
    /// @notice Lock tokens in the vault
    /// @param _from Address from which tokens are locked
    /// @param _amount Amount of tokens to lock
    function lockTokenInVault(address _from, uint256 _amount) external;

    /// @notice Release tokens from the vault
    /// @param _to Address to which tokens are released
    /// @param _amount Amount of tokens to release
    function releaseTokenInVault(address _to, uint256 _amount) external;
}

/// @title Messenger Contract
/// @notice A contract for sending and receiving string data across chains using Chainlink CCIP
contract CrossChain_Messanger is CCIPReceiver, OwnerIsCreator {
    using SafeERC20 for IERC20;

    /// @notice Struct to store client data
    struct ClientData {
        uint256 amount;
        bool exists;
    }

    // Custom errors
    error NotEnoughBalance(uint256 currentBalance, uint256 calculatedFees);
    error NothingToWithdraw();
    error FailedToWithdrawEth(address owner, address target, uint256 value);
    error DestinationChainNotAllowlisted(uint64 destinationChainSelector);
    error SourceChainNotAllowlisted(uint64 sourceChainSelector);
    error SenderNotAllowlisted(address sender);
    error InvalidReceiverAddress();

    /// @notice Emitted when a message is sent to another chain
    event MessageSent(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address receiver,
        string text,
        uint256 amount,
        address client,
        address feeToken,
        uint256 fees
    );

    /// @notice Emitted when a message is received from another chain
    event MessageReceived(
        bytes32 indexed messageId,
        uint64 indexed sourceChainSelector,
        address sender,
        string text,
        uint256 amount,
        address client
    );

    bytes32 private s_lastReceivedMessageId;
    string private s_lastReceivedText;
    address[] public clientAddresses;

    mapping(address => ClientData) public clientDataMap;
    mapping(uint64 => bool) public allowlistedDestinationChains;
    mapping(uint64 => bool) public allowlistedSourceChains;
    mapping(address => bool) public allowlistedSenders;

    IERC20 private s_linkToken;
    CCIP_TokenVault private vault;

    /// @notice Constructor
    /// @param _router The address of the router contract
    /// @param _link The address of the LINK token contract
    /// @param _vault The address of the vault contract
    constructor(address _router, address _link, address _vault) CCIPReceiver(_router) {
        s_linkToken = IERC20(_link);
        vault = CCIP_TokenVault(_vault);
    }

    /// @notice Modifier to check if the destination chain is allowlisted
    modifier onlyAllowlistedDestinationChain(uint64 _destinationChainSelector) {
        if (!allowlistedDestinationChains[_destinationChainSelector])
            revert DestinationChainNotAllowlisted(_destinationChainSelector);
        _;
    }

    /// @notice Modifier to check if the source chain and sender are allowlisted
    modifier onlyAllowlisted(uint64 _sourceChainSelector, address _sender) {
        if (!allowlistedSourceChains[_sourceChainSelector])
            revert SourceChainNotAllowlisted(_sourceChainSelector);
        if (!allowlistedSenders[_sender]) revert SenderNotAllowlisted(_sender);
        _;
    }

    /// @notice Modifier to validate the receiver address
    modifier validateReceiver(address _receiver) {
        if (_receiver == address(0)) revert InvalidReceiverAddress();
        _;
    }

    /// @notice Allowlist a destination chain
    /// @param _destinationChainSelector The selector of the destination chain
    /// @param allowed Whether the chain should be allowlisted
    function allowlistDestinationChain(
        uint64 _destinationChainSelector,
        bool allowed
    ) external onlyOwner {
        allowlistedDestinationChains[_destinationChainSelector] = allowed;
    }

    /// @notice Allowlist a source chain
    /// @param _sourceChainSelector The selector of the source chain
    /// @param allowed Whether the chain should be allowlisted
    function allowlistSourceChain(
        uint64 _sourceChainSelector,
        bool allowed
    ) external onlyOwner {
        allowlistedSourceChains[_sourceChainSelector] = allowed;
    }

    /// @notice Allowlist a sender
    /// @param _sender The address of the sender
    /// @param allowed Whether the sender should be allowlisted
    function allowlistSender(address _sender, bool allowed) external onlyOwner {
        allowlistedSenders[_sender] = allowed;
    }

    /// @notice Send a message to another chain, paying fees in LINK
    /// @param _destinationChainSelector The selector of the destination chain
    /// @param _receiver The address of the receiver on the destination chain
    /// @param _text The text message to send
    /// @param _amount The amount of tokens to transfer
    /// @param _client The address of the client
    /// @return messageId The ID of the sent message
    function sendMessagePayLINK(
        uint64 _destinationChainSelector,
        address _receiver,
        string calldata _text,
        uint256 _amount,
        address _client
    )
        external
        onlyOwner
        onlyAllowlistedDestinationChain(_destinationChainSelector)
        validateReceiver(_receiver)
        returns (bytes32 messageId)
    {
        Client.EVM2AnyMessage memory evm2AnyMessage = _buildCCIPMessage(
            _receiver,
            _text,
            _amount,
            _client,
            address(s_linkToken)
        );

        IRouterClient router = IRouterClient(this.getRouter());

        uint256 fees = router.getFee(_destinationChainSelector, evm2AnyMessage);

        if (fees > s_linkToken.balanceOf(address(this)))
            revert NotEnoughBalance(s_linkToken.balanceOf(address(this)), fees);

        s_linkToken.approve(address(router), fees);

        messageId = router.ccipSend(_destinationChainSelector, evm2AnyMessage);

        vault.lockTokenInVault(msg.sender, _amount);

        emit MessageSent(
            messageId,
            _destinationChainSelector,
            _receiver,
            _text,
            _amount,
            _client,
            address(s_linkToken),
            fees
        );

        return messageId;
    }

    /// @notice Handle a received message
    /// @param any2EvmMessage The received message
    function _ccipReceive(
        Client.Any2EVMMessage memory any2EvmMessage
    )
        internal
        override
        onlyAllowlisted(
            any2EvmMessage.sourceChainSelector,
            abi.decode(any2EvmMessage.sender, (address))
        )
    {
        s_lastReceivedMessageId = any2EvmMessage.messageId;
        s_lastReceivedText = abi.decode(any2EvmMessage.data, (string));
        (string memory text, uint256 amount, address client) = abi.decode(any2EvmMessage.data, (string, uint256, address));

        if (!clientDataMap[client].exists) {
            clientAddresses.push(client);
            clientDataMap[client].exists = true;
        }
        clientDataMap[client].amount = amount;
        vault.releaseTokenInVault(client, amount);

        emit MessageReceived(
            any2EvmMessage.messageId,
            any2EvmMessage.sourceChainSelector,
            abi.decode(any2EvmMessage.sender, (address)),
            text,
            amount,
            client
        );
    }

    /// @notice Build a CCIP message
    /// @param _receiver The address of the receiver
    /// @param _text The text message
    /// @param _amount The amount of tokens
    /// @param _client The address of the client
    /// @param _feeTokenAddress The address of the fee token
    /// @return An EVM2AnyMessage struct
    function _buildCCIPMessage(
        address _receiver,
        string calldata _text,
        uint256 _amount,
        address _client,
        address _feeTokenAddress
    ) private pure returns (Client.EVM2AnyMessage memory) {
        return
            Client.EVM2AnyMessage({
                receiver: abi.encode(_receiver),
                data: abi.encode(_text, _amount, _client),
                tokenAmounts: new Client.EVMTokenAmount[](0),
                extraArgs: Client._argsToBytes(
                    Client.EVMExtraArgsV1({gasLimit: 400_000})
                ),
                feeToken: _feeTokenAddress
            });
    }

    /// @notice Get details of the last received message
    /// @return messageId The ID of the last received message
    /// @return text The text of the last received message
    function getLastReceivedMessageDetails()
        external
        view
        returns (bytes32 messageId, string memory text)
    {
        return (s_lastReceivedMessageId, s_lastReceivedText);
    }

    /// @notice Receive function to accept Ether
    receive() external payable {}

    /// @notice Withdraw Ether from the contract
    /// @param _beneficiary The address to receive the withdrawn Ether
    function withdraw(address _beneficiary) public onlyOwner {
        uint256 amount = address(this).balance;
        if (amount == 0) revert NothingToWithdraw();
        (bool sent, ) = _beneficiary.call{value: amount}("");
        if (!sent) revert FailedToWithdrawEth(msg.sender, _beneficiary, amount);
    }

    /// @notice Withdraw ERC20 tokens from the contract
    /// @param _beneficiary The address to receive the withdrawn tokens
    /// @param _token The address of the ERC20 token to withdraw
    function withdrawToken(
        address _beneficiary,
        address _token
    ) public onlyOwner {
        uint256 amount = IERC20(_token).balanceOf(address(this));
        if (amount == 0) revert NothingToWithdraw();
        IERC20(_token).safeTransfer(_beneficiary, amount);
    }

    /// @notice Get all client data
    /// @return An array of client addresses and an array of their corresponding amounts
    function getAllClientData() external view returns (address[] memory, uint256[] memory) {
        uint256[] memory amounts = new uint256[](clientAddresses.length);
        
        for (uint i = 0; i < clientAddresses.length; i++) {
            amounts[i] = clientDataMap[clientAddresses[i]].amount;
        }

        return (clientAddresses, amounts);
    }
}