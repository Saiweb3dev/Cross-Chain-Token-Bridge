// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouterClient} from "@chainlink/contracts-ccip/src/v0.8/ccip/interfaces/IRouterClient.sol";
import {OwnerIsCreator} from "@chainlink/contracts-ccip/src/v0.8/shared/access/OwnerIsCreator.sol";
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol";
import {IERC20} from "@chainlink/contracts-ccip/src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@chainlink/contracts-ccip/src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title ICCIPvault
 * @dev Interface for interacting with the CCIPvault contract
 */
interface ICCIPvault {
    /**
     * @dev Lock tokens from a sender
     * @param _from Sender address
     * @param _amount Amount to lock
     * @return bool Whether the operation was successful
     */
    function lockTokens(address _from, uint256 _amount) external returns (bool);

    /**
     * @dev Release locked tokens to a receiver
     * @param _to Receiver address
     * @param _amount Amount to release
     * @return bool Whether the operation was successful
     */
    function releaseTokens(address _to, uint256 _amount) external returns (bool);
}

contract Cross_Chain_Messenger is CCIPReceiver, OwnerIsCreator {
    using SafeERC20 for IERC20;

    // Custom error messages for better error handling
    error NotEnoughBalance(uint256 currentBalance, uint256 calculatedFees);
    error NothingToWithdraw();
    error FailedToWithdrawEth(address owner, address target, uint256 value);
    error DestinationChainNotAllowlisted(uint64 destinationChainSelector);
    error SourceChainNotAllowlisted(uint64 sourceChainSelector);
    error SenderNotAllowlisted(address sender);
    error InvalidReceiverAddress();
    error FailedToLockTokens(address from, uint256 amount);
    error FailedToReleaseTokens(address to, uint256 amount);

    // Events emitted during message operations
    event MessageSent(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address receiver,
        string text,
        address feeToken,
        uint256 fees
    );
    event MessageReceived(
        bytes32 indexed messageId,
        uint64 indexed sourceChainSelector,
        address sender,
        string text
    );
    event TokensLocked(address indexed from, uint256 amount);
    event TokensReleased(address indexed to, uint256 amount);

    // State variables
    bytes32 private s_lastReceivedMessageId;
    string private s_lastReceivedText;

    // Allowlist mappings
    mapping(uint64 => bool) public allowlistedDestinationChains;
    mapping(uint64 => bool) public allowlistedSourceChains;
    mapping(address => bool) public allowlistedSenders;

    // Contract dependencies
    ICCIPvault public vault;
    IERC20 private s_linkToken;

    constructor(
        address _router,
        address _link,
        address _vault
    ) CCIPReceiver(_router) {
        s_linkToken = IERC20(_link);
        vault = ICCIPvault(_vault);
    }

    /**
     * @dev Modifier to ensure the destination chain is allowlisted
     * @param _destinationChainSelector Chain selector to validate
     */
    modifier onlyAllowlistedDestinationChain(uint64 _destinationChainSelector) {
        if (!allowlistedDestinationChains[_destinationChainSelector])
            revert DestinationChainNotAllowlisted(_destinationChainSelector);
        _;
    }

    /**
     * @dev Modifier to ensure both source chain and sender are allowlisted
     * @param _sourceChainSelector Source chain selector
     * @param _sender Sender address
     */
    modifier onlyAllowlisted(uint64 _sourceChainSelector, address _sender) {
        if (!allowlistedSourceChains[_sourceChainSelector])
            revert SourceChainNotAllowlisted(_sourceChainSelector);
        if (!allowlistedSenders[_sender]) revert SenderNotAllowlisted(_sender);
        _;
    }

    /**
     * @dev Modifier to validate the receiver address
     * @param _receiver Address to validate
     */
    modifier validateReceiver(address _receiver) {
        if (_receiver == address(0)) revert InvalidReceiverAddress();
        _;
    }

    /**
     * @dev Allowlist a destination chain
     * @param _destinationChainSelector Chain selector
     * @param _allowed Whether to allowlist
     */
    function allowlistDestinationChain(
        uint64 _destinationChainSelector,
        bool allowed
    ) external onlyOwner {
        allowlistedDestinationChains[_destinationChainSelector] = allowed;
    }

    /**
     * @dev Allowlist a source chain
     * @param _sourceChainSelector Chain selector
     * @param _allowed Whether to allowlist
     */
    function allowlistSourceChain(uint64 _sourceChainSelector, bool allowed)
        external
        onlyOwner
    {
        allowlistedSourceChains[_sourceChainSelector] = allowed;
    }

    /**
     * @dev Allowlist a sender
     * @param _sender Sender address
     * @param _allowed Whether to allowlist
     */
    function allowlistSender(address _sender, bool allowed) external onlyOwner {
        allowlistedSenders[_sender] = allowed;
    }

    /**
     * @dev Send a message using LINK token
     * @param _destinationChainSelector Destination chain selector
     * @param _receiver Receiver address
     * @param _text Message text
     * @return bytes32 Message ID
     */
    function sendMessagePayLINK(
        uint64 _destinationChainSelector,
        address _receiver,
        string calldata _text
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
            address(s_linkToken)
        );

        IRouterClient router = IRouterClient(this.getRouter());

        uint256 amount = stringToUint(_text);

        bool result = vault.lockTokens(msg.sender, amount);
        if (!result) revert FailedToLockTokens(msg.sender, amount);
        emit TokensLocked(msg.sender, amount);

        uint256 fees = router.getFee(_destinationChainSelector, evm2AnyMessage);

        if (fees > s_linkToken.balanceOf(address(this)))
            revert NotEnoughBalance(s_linkToken.balanceOf(address(this)), fees);

        s_linkToken.approve(address(router), fees);

        messageId = router.ccipSend(_destinationChainSelector, evm2AnyMessage);

        emit MessageSent(
            messageId,
            _destinationChainSelector,
            _receiver,
            _text,
            address(s_linkToken),
            fees
        );

        return messageId;
    }

    /**
     * @dev Receive a message from any chain
     * @param any2EvmMessage Received message details
     */
    function _ccipReceive(Client.Any2EVMMessage memory any2EvmMessage)
        internal
        override
        onlyAllowlisted(
            any2EvmMessage.sourceChainSelector,
            abi.decode(any2EvmMessage.sender, (address))
        )
    {
        s_lastReceivedMessageId = any2EvmMessage.messageId;
        s_lastReceivedText = abi.decode(any2EvmMessage.data, (string));

        uint256 amount = stringToUint(s_lastReceivedText);
        address sender = abi.decode(any2EvmMessage.sender, (address));

        bool result = vault.releaseTokens(sender, amount);
        if (!result) revert FailedToReleaseTokens(sender, amount);

        emit TokensReleased(sender, amount);

        emit MessageReceived(
            any2EvmMessage.messageId,
            any2EvmMessage.sourceChainSelector,
            sender,
            s_lastReceivedText
        );
    }

    /**
     * @dev Build a CCIP message
     * @param _receiver Receiver address
     * @param _text Message text
     * @param _feeTokenAddress Fee token address
     * @return Client.EVM2AnyMessage Built message
     */
    function _buildCCIPMessage(
        address _receiver,
        string calldata _text,
        address _feeTokenAddress
    ) private pure returns (Client.EVM2AnyMessage memory) {
        return
            Client.EVM2AnyMessage({
                receiver: abi.encode(_receiver),
                data: abi.encode(_text),
                tokenAmounts: new Client.EVMTokenAmount[](0),
                extraArgs: Client._argsToBytes(
                    Client.EVMExtraArgsV1({gasLimit: 200_000})
                ),
                feeToken: _feeTokenAddress
            });
    }

    /**
     * @dev Get the last received message details
     * @return bytes32 Last received message ID
     * @return string Last received message text
     */
    function getLastReceivedMessageDetails()
        external
        view
        returns (bytes32 messageId, string memory text)
    {
        return (s_lastReceivedMessageId, s_lastReceivedText);
    }

    /**
     * @dev Receive Ether sent to this contract
     */
    receive() external payable {}

    /**
     * @dev Withdraw ETH to a beneficiary address
     * @param _beneficiary Beneficiary address
     */
    function withdraw(address _beneficiary) public onlyOwner {
        uint256 amount = address(this).balance;

        if (amount == 0) revert NothingToWithdraw();

        (bool sent, ) = _beneficiary.call{value: amount}("");

        if (!sent) revert FailedToWithdrawEth(msg.sender, _beneficiary, amount);
    }

    /**
     * @dev Withdraw a specific token to a beneficiary address
     * @param _beneficiary Beneficiary address
     * @param _token Token address
     */
    function withdrawToken(address _beneficiary, address _token)
        public
        onlyOwner
    {
        uint256 amount = IERC20(_token).balanceOf(address(this));

        if (amount == 0) revert NothingToWithdraw();

        IERC20(_token).safeTransfer(_beneficiary, amount);
    }

    /**
     * @dev Convert a string to a uint256
     * @param str String to convert
     * @return uint256 Converted value
     */
    function stringToUint(string memory str) internal pure returns (uint256) {
        bytes memory b = bytes(str);
        uint256 result = 0;
        for (uint256 i = 0; i < b.length; i++) {
            uint8 c = uint8(b[i]);
            if (c >= 48 && c <= 57) {
                result = result * 10 + (c - 48);
            }
        }
        return result;
    }
}
