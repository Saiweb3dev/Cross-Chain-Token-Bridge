import * as React from 'react';
import { useAccount, useDisconnect } from 'wagmi';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

/**
 * Account component displays the connected wallet address and allows the user to disconnect their wallet.
 */
export function Account() {
  const { address } = useAccount(); // Hook to get the connected account's address
  const { disconnect } = useDisconnect(); // Hook to handle wallet disconnection
  const [showModal, setShowModal] = React.useState(false); // State to control the visibility of the disconnect modal

  /**
   * Helper function to format the wallet address for display.
   * Shows the first 4 characters, ellipsis (...), and the last 4 characters.
   * @param {string} address - The wallet address to format.
   * @returns {string} Formatted address string.
   */
  const formatAddress = (address:string) => `${address.substring(0, 4)}...${address.substring(address.length - 4)}`;

  /**
   * Inline functional component for rendering the disconnect modal.
   * @param {{ isOpen: boolean; onClose: () => void; children: React.ReactNode }} props - Props passed to the Modal component.
   * @returns JSX.Element - The modal element.
   */
  const Modal = ({ isOpen, onClose, children }: ModalProps) => {
    if (!isOpen) return null;

    return (
      <div className="fixed inset-0 bg-purple-400 bg-opacity-50 flex items-center justify-center z-50 ">
        <div className="bg-purple-700 p-4 rounded shadow-lg w-full max-w-xs space-y-4 border-2 border-white">
          {children}
          <button 
            className='bg-white text-black px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-200 transition duration-300 w-full'
            onClick={onClose}
          >
            Close
          </button>
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col items-center space-y-4 mr-6">
      {/* Display the formatted address as a button if an address exists */}
      {address && (
        <button 
          className="text-sm font-medium text-gray-300 hover:text-gray-500 border px-4 rounded-lg"
          onClick={() => setShowModal(true)} // Show modal when address button is clicked
        >
          {formatAddress(address)}
        </button>
      )}

      {/* Render the disconnect modal inline */}
      <Modal isOpen={showModal} onClose={() => setShowModal(false)}>
        <h2 className="text-xl text-white font-semibold">Disconnect Wallet</h2>
        <p className="text-white">Are you sure you want to disconnect?</p>
        <button 
          className="bg-purple-700 text-white px-4 py-2 rounded-md text-sm font-medium border hover:bg-purple-500 transition duration-300 w-full mt-4"
          onClick={() => {
            disconnect(); // Disconnect wallet
            setShowModal(false); // Close modal after disconnecting
          }}
        >
          Disconnect
        </button>
      </Modal>
    </div>
  );
}