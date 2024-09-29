import { useState, useEffect, useCallback } from 'react';
import Web3 from 'web3';
import { Contract } from 'web3-eth-contract';
import { AbiItem } from 'web3-utils';

interface UseSmartContractProps {
  address: string;
  abi: any;
  chainId: number | string;
}
interface UseSmartContractReturn {
  contract: Contract<AbiItem[]> | null;
  account: string | null;
  loading: boolean;
  error: Error | null;
  callMethod: (method: string, ...args: any[]) => Promise<any>;
  sendTransaction: (method: string, ...args: any[]) => Promise<any>;
}

const transformAbi = (abi: any): AbiItem[] => {
  const transformedAbi: AbiItem[] = [];

  // Transform constructor
  if (abi.Constructor) {
    transformedAbi.push({
      type: 'constructor',
      inputs: abi.Constructor.Inputs.map((input: any) => ({
        name: input.Name,
        type: getTypeString(input.Type)
      })),
      stateMutability: abi.Constructor.StateMutability
    });
  }

  // Transform methods
  Object.values(abi.Methods).forEach((method: any) => {
    transformedAbi.push({
      type: 'function',
      name: method.Name,
      inputs: method.Inputs.map((input: any) => ({
        name: input.Name,
        type: getTypeString(input.Type)
      })),
      outputs: method.Outputs.map((output: any) => ({
        name: output.Name,
        type: getTypeString(output.Type)
      })),
      stateMutability: method.StateMutability,
      constant: method.Constant,
      payable: method.Payable
    });
  });

  // Transform events (if needed)
  // ... (add event transformation logic here if required)

  return transformedAbi;
};

const getTypeString = (type: any): string => {
  switch (type.T) {
    case 1: return `uint${type.Size}`;
    case 2: return 'bool';
    case 3: return 'string';
    case 7: return 'address';
    // Add more cases as needed
    default: return 'unknown';
  }
};

export const useSmartContract = ({ address, abi,chainId }: UseSmartContractProps): UseSmartContractReturn => {
  const [web3, setWeb3] = useState<Web3 | null>(null);
  const [contract, setContract] = useState<Contract<typeof abi> | null>(null);
  const [account, setAccount] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const initWeb3 = async () => {
      try {
        if (typeof window !== 'undefined' && typeof window.ethereum !== 'undefined') {
          const web3Instance = new Web3(window.ethereum);
          setWeb3(web3Instance);

          const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
          setAccount(accounts[0]);

          // Check if the current network matches the expected chainId
          const currentChainId = await web3Instance.eth.getChainId();
          if (currentChainId.toString() !== chainId.toString()) {
            throw new Error(`Please switch to the correct network. Expected chainId: ${chainId}, Current chainId: ${currentChainId}`);
          }

          const transformedAbi = transformAbi(abi);
          const contractInstance = new web3Instance.eth.Contract(transformedAbi, address);
          setContract(contractInstance);
        } else {
          throw new Error('Please install MetaMask to use this dApp');
        }
      } catch (err) {
        console.error('Error in initWeb3:', err);
        setError(err instanceof Error ? err : new Error('An unknown error occurred'));
      } finally {
        setLoading(false);
      }
    };

    initWeb3();
  }, [abi, address, chainId]);

  const callMethod = useCallback(async (method: string, ...args: any[]): Promise<any> => {
    if (!contract) {
      throw new Error('Contract is not initialized');
    }
    try {
      const contractMethod = contract.methods[method];
      if (typeof contractMethod !== 'function') {
        throw new Error(`Method ${method} is not a function on this contract`);
      }
      return await contractMethod(...args).call({ from: account || undefined });
    } catch (err) {
      setError(err instanceof Error ? err : new Error('An error occurred while calling the method'));
      throw err;
    }
  }, [contract, account]);

  const sendTransaction = useCallback(async (
    method: string, 
    args: any[], 
    options: { gasLimit?: number } = {}
  ): Promise<any> => {
    if (!contract || !account) {
      throw new Error('Contract or account is not initialized');
    }
  
    try {
      console.log('Sending transaction:', method, args);
      const contractMethod = contract.methods[method];
      if (typeof contractMethod !== 'function') {
        throw new Error(`Method ${method} is not a function on this contract`);
      }

      const defaultGasLimit = 300000; // Increased default gas limit
      const gasLimit = options.gasLimit || defaultGasLimit;

      let gasEstimate;
      try {
        gasEstimate = await contractMethod(...args).estimateGas({ from: account });
        console.log('Estimated gas:', gasEstimate);
      } catch (estimateError) {
        console.warn('Gas estimation failed:', estimateError);
        console.log('Using default gas limit:', gasLimit);
        gasEstimate = gasLimit;
      }

      return await contractMethod(...args).send({ 
        from: account, 
        gas: Math.min(Number(gasEstimate) * 1.2, Number(gasLimit)).toString() 
      });
    } catch (err) {
      console.error('Error in sendTransaction:', err);
      setError(err instanceof Error ? err : new Error('An error occurred while sending the transaction'));
      throw err;
    }
  }, [contract, account]);

  return { contract, account, loading, error, callMethod, sendTransaction };
};