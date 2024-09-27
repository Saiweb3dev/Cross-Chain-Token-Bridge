"use client"
import React, { createContext, useState, useContext, ReactNode } from 'react';

type ChainContextType = {
  chainId: string;
  setChainId: (chainId: string) => void;
};

const ChainContext = createContext<ChainContextType | undefined>(undefined);

interface ChainProviderProps {
  children: ReactNode;
}

export const ChainProvider: React.FC<ChainProviderProps> = ({ children }) => {
  const [chainId, setChainId] = useState('11155111'); // Default to Sepolia

  return (
    <ChainContext.Provider value={{ chainId, setChainId }}>
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