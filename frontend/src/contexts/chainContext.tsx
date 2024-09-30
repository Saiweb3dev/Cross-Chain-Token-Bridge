"use client"
import React, { createContext, useState, useContext, ReactNode, useEffect } from 'react';
import { ethers } from 'ethers';

console.log('ChainContext is being loaded');

type ChainContextType = {
  chainId: string;
  setChainId: (chainId: string) => void;
  isNetworkSwitchRequired: boolean;
  switchNetwork: (chainId: string) => Promise<void>;
};

export const ChainContext = createContext<ChainContextType | undefined>(undefined);

interface ChainProviderProps {
  children: ReactNode;
}

 export const ChainProvider: React.FC<ChainProviderProps> = ({ children }) => {
  const [selectedChainId, setSelectedChainId] = useState('11155111'); // Default to Sepolia
  const [currentChainId, setCurrentChainId] = useState<string | null>(null);

  const isNetworkSwitchRequired = currentChainId !== selectedChainId;

  useEffect(() => {
    const checkNetwork = async () => {
      if (window.ethereum) {
        const provider = new ethers.BrowserProvider(window.ethereum);
        const network = await provider.getNetwork();
        setCurrentChainId(network.chainId.toString());
      }
    };

    checkNetwork();

    if (window.ethereum) {
      window.ethereum.on('chainChanged', (chainId: string) => {
        setCurrentChainId(parseInt(chainId).toString());
      });
    }

    return () => {
      if (window.ethereum) {
        window.ethereum.removeAllListeners('chainChanged');
      }
    };
  }, []);

  const switchNetwork = async (chainId: string) => {
    if (window.ethereum) {
      try {
        await window.ethereum.request({
          method: 'wallet_switchEthereumChain',
          params: [{ chainId: `0x${parseInt(chainId).toString(16)}` }],
        });
      } catch (switchError: any) {
        // This error code indicates that the chain has not been added to MetaMask.
        if (switchError.code === 4902) {
          // You could add logic here to add the chain to the user's wallet
          console.log('This network is not available in your metamask, please add it')
        }
        console.log('Failed to switch to the network')
      }
    }
  };

  const setChainId = (chainId: string) => {
    setSelectedChainId(chainId);
    switchNetwork(chainId);
  };

  return (
    <ChainContext.Provider value={{ 
      chainId: selectedChainId, 
      setChainId, 
      isNetworkSwitchRequired,
      switchNetwork
    }}>
      {children}
    </ChainContext.Provider>
  );
};

 export const useChain = (): ChainContextType => {
  const context = useContext(ChainContext);
  if (context === undefined) {
    throw new Error('useChain must be used within a ChainProvider');
  }
  return context;
};

