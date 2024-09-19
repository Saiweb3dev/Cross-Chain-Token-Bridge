import * as React from 'react';
import { Connector, useConnect } from 'wagmi'; // wagmi hooks and types for wallet connection

/**
 * WalletOptions component renders a button to open a modal with available wallet connection options.
 * Users can select a wallet to connect to the dApp.
 */
export function WalletOptions() {
  const { connectors, connect } = useConnect(); // wagmi hook to get available connectors and connect function
  const [showModal, setShowModal] = React.useState(false); // State to control the visibility of the wallet options modal

  /**
   * Function to show the modal when the "Connect" button is clicked.
   */
  const handleConnectClick = () => {
    setShowModal(true); // Show the modal when "Connect" is clicked
  };

  /**
   * Function to handle wallet selection.
   * Connects to the selected wallet and closes the modal.
   * @param {Connector} connector - The selected wallet connector.
   */
  const handleWalletSelect = (connector:any) => {
    connect({ connector }); // Connect to the selected wallet
    setShowModal(false); // Close the modal after a wallet is selected
  };

  return (
    <div className="flex flex-col items-center justify-center space-y-4 mr-6">
      {/* Button to open the wallet options modal */}
      {!showModal && (
        <button 
          onClick={handleConnectClick}
          className="bg-purple-700 border-2 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-purple-500"
        >
          Connect
        </button>
      )}

      {/* Modal displaying available wallet connection options */}
      {showModal && (
        <div className="fixed inset-0 p-16 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-purple-600 text-white text-center px-6 py-28 rounded-md text-sm font-medium duration-300 hover:bg-purple-500 border-2">
            <h2 className="text-4xl font-semibold">Connect Wallet</h2>
            <div className='flex flex-row space-x-6 m-6'>
              {/* Map over available connectors to render buttons for each wallet option */}
              {connectors.map((connector) => (
                <button 
                  key={connector.uid} 
                  onClick={() => handleWalletSelect(connector)}
                  className="bg-purple-200 text-black px-6 py-2 rounded hover:bg-purple-300 hover:text-gray-700 transition duration-300 w-fit "
                >
                  {connector.name}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}