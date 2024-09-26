import { useState, useCallback, useEffect } from 'react';
import { ethers } from 'ethers';

interface UseContractFunctionProps {
  abi: ethers.InterfaceAbi;
  address: string;
  functionName: string;
  args?: any[];
}

export const useContractFunction = ({
  abi,
  address,
  functionName,
  args = [],
}: UseContractFunctionProps) => {
  const [contract, setContract] = useState<ethers.Contract | null>(null);
  const [signer, setSigner] = useState<ethers.Signer | null>(null);
  const [provider, setProvider] = useState<ethers.BrowserProvider | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    const initializeEthers = async () => {
      if (typeof window.ethereum !== 'undefined') {
        const provider = new ethers.BrowserProvider(window.ethereum);
        setProvider(provider);
        try {
          const signer = await provider.getSigner();
          setSigner(signer);
          const contract = new ethers.Contract(address, abi, signer);
          setContract(contract);
        } catch (err) {
          setError(new Error('Failed to get signer. Please connect to MetaMask.'));
        }
      } else {
        setError(new Error('Please install MetaMask!'));
      }
    };

    initializeEthers();
  }, [abi, address]);

  const execute = useCallback(async (...callArgs: any[]) => {
    if (!contract) {
      throw new Error('Contract is not initialized');
    }

    setIsLoading(true);
    setError(null);
    setData(null);

    try {
      const contractFunction = contract[functionName as keyof typeof contract];
      if (typeof contractFunction !== 'function') {
        throw new Error(`Function ${functionName} not found on contract`);
      }

      const result = await contractFunction(...(callArgs.length > 0 ? callArgs : args));
      
      // If the function is a write function, wait for the transaction to be mined
      if (result && typeof result.wait === 'function') {
        const receipt = await result.wait();
        setData(receipt);
      } else {
        setData(result);
      }

      return result;
    } catch (err) {
      setError(err instanceof Error ? err : new Error('An unknown error occurred'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [contract, functionName, args]);

  const isReadFunction = useCallback(() => {
    if (!contract) return false;
    const functionFragment = contract.interface.getFunction(functionName);
    return functionFragment?.constant || functionFragment?.stateMutability === 'view' || functionFragment?.stateMutability === 'pure';
  }, [contract, functionName]);

  return {
    execute,
    data,
    isLoading,
    error,
    isReadFunction: isReadFunction(),
    contract,
    signer,
    provider,
  };
};