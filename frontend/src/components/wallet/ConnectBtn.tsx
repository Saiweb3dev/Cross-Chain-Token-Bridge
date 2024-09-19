// Import necessary dependencies and components
"use client"
import React, { useEffect, useState } from 'react';
import { useAccount } from 'wagmi'; // wagmi hook to check if a wallet is connected
import { Account } from './Account'; // Component to display account details
import {WalletOptions} from './WalletOptions'; // Component to offer wallet connection options

/**
 * ConnectWallet component decides whether to show the Account component (if a wallet is connected)
 * or the WalletOptions component (if no wallet is connected).
 */
function ConnectBtn() {
  const { isConnected } = useAccount(); // wagmi hook to check the connection status of the wallet
  const [isClient, setIsClient] = useState(false); // State to track if component is rendered on the client-side

  /**
   * useEffect hook to set isClient state to true once the component has mounted.
   * This ensures that the component does not attempt to render based on isConnected
   * until it is certain that the rendering is happening on the client-side.
   */
  useEffect(() => {
    setIsClient(true); // Set to true once component has mounted
  }, []);

  /**
   * Conditional rendering based on the connection status and whether the component
   * is being rendered on the client-side.
   */
  // If not yet rendered on the client-side, render nothing
  if (!isClient) return null;

  // If a wallet is connected, render the Account component
  if (isConnected) return <Account />;

  // If no wallet is connected, render the WalletOptions component
  return <WalletOptions />;
}

export default ConnectBtn;