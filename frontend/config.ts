import { http, createConfig } from 'wagmi'; // Import necessary functions from wagmi
import { base, mainnet } from 'wagmi/chains'; // Import chain definitions
import { injected, metaMask, safe } from 'wagmi/connectors'; // Import connector definitions

const projectId = '<WALLETCONNECT_PROJECT_ID>'; // Placeholder for WalletConnect project ID

/**
 * Creates and exports a wagmi configuration object.
 * This configuration specifies the chains to support, the connectors for wallet connections,
 * and the HTTP transport method for interacting with the blockchain.
 */
export const config = createConfig({
  /**
   * Chains supported by the dApp.
   * Includes Ethereum mainnet and a base chain (e.g., a Layer 2 or testnet).
   */
  chains: [mainnet, base],

  /**
   * Connectors define how users can connect their wallets to the dApp.
   * This configuration includes support for injected wallets (like MetaMask),
   * MetaMask specifically, Safe (Gnosis Safe), and WalletConnect.
   * Note: WalletConnect is commented out for demonstration purposes.
   */
  connectors: [
    injected(), // Supports wallets injected into the page (e.g., MetaMask)
    // walletConnect({ projectId }), // Uncomment to enable WalletConnect support
    metaMask({
      /**
       * Metadata for the dApp when connecting through MetaMask.
       * This information is displayed to the user during connection.
       */
      dappMetadata: {
        name: "KITE_EDUCATIONAL_DAPP", // Name of the dApp
        url: "http://localhost:3000", // URL of the dApp
      },
    }),
    safe(), // Supports Gnosis Safe connections
  ],

  /**
   * Transports specify how to interact with the blockchain for different chains.
   * Here, HTTP transport is configured for both mainnet and base chain.
   */
  transports: {
    [mainnet.id]: http(), // HTTP transport for Ethereum mainnet
    [base.id]: http(), // HTTP transport for the base chain
  },
});