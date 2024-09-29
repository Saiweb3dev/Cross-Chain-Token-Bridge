import React, { useState,useEffect } from 'react';
import { motion } from 'framer-motion';
import { ArrowRight } from 'lucide-react';
import { useSmartContract } from '../../../hooks/useContractFunction';
import { AbiItem } from 'web3-utils';
import Web3 from 'web3';
import { useChain } from '@/contexts/chainContext';

interface ContractFunctionProps {
  title: string;
  contractAddress: string;
  abi: AbiItem | AbiItem[];
  functionName: string;
}

export default function ContractFunction({ title, contractAddress, abi, functionName }: ContractFunctionProps) {
  const [tokenAmount, setTokenAmount] = useState('');
  const [recipientAddress, setRecipientAddress] = useState('');
  const [isHovered, setIsHovered] = useState(false);
  const [isSelfMint, setIsSelfMint] = useState(true);
  const [loading, setLoading] = useState(false);

  const { chainId } = useChain();

  // Ensure abi is always an array
  const { contract, account, sendTransaction, error } = useSmartContract({ address: contractAddress, abi,chainId });

  useEffect(() => {
    if (contract) {
      console.log('Available methods:', Object.keys(contract.methods));
    }
  }, [contract]);

  const handleAction = async () => {
    if (loading) return;
    setLoading(true);
    try {
      const web3 = new Web3(window.ethereum);
      const amount = web3.utils.toWei(tokenAmount, 'ether');
      const recipient = isSelfMint ? account : recipientAddress;
  
      console.log('Calling function:', functionName);
      console.log('Arguments:', recipient, amount);
  
      const result = await sendTransaction(functionName, [recipient, amount], { gasLimit: 300000 });
      console.log('Transaction result:', result);
  
      console.log(`${functionName} successful`);
      // Add success feedback here
    } catch (err) {
      console.error(`Error in ${functionName}:`, err);
      if (err instanceof Error) {
        // Add more detailed error handling here
        if (err.message.includes('Internal JSON-RPC error')) {
          console.error('Possible contract execution error. Check the contract code and parameters.');
        }
      }
      // Add error handling UI here
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <motion.div 
      className="bg-white rounded-xl shadow-lg overflow-hidden w-full max-w-sm mx-auto"
      initial={{ y: 20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ type: "spring", stiffness: 260, damping: 20 }}
      onHoverStart={() => setIsHovered(true)}
      onHoverEnd={() => setIsHovered(false)}
    >
      <div className="p-6">
        <motion.h2 
          className="text-xl font-bold text-purple-800 mb-4"
          initial={{ y: -20 }}
          animate={{ y: 0 }}
          transition={{ type: "spring", stiffness: 300, damping: 10 }}
        >
          {title}
        </motion.h2>
        <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
          <div>
            <label htmlFor={`tokenAmount-${title}`} className="block text-sm font-medium text-gray-700 mb-1">
              Token Amount
            </label>
            <input
              type="number"
              id={`tokenAmount-${title}`}
              value={tokenAmount}
              onChange={(e) => setTokenAmount(e.target.value)}
              className="block text-purple-300 w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50 transition duration-200"
              placeholder="Enter amount"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Recipient
            </label>
            <div className="flex items-center space-x-2">
              <button
                type="button"
                onClick={() => setIsSelfMint(true)}
                className={`px-3 py-1 rounded-md text-sm ${isSelfMint ? 'bg-purple-600 text-white' : 'bg-gray-200 text-gray-700'}`}
              >
                Self
              </button>
              <button
                type="button"
                onClick={() => setIsSelfMint(false)}
                className={`px-3 py-1 rounded-md text-sm ${!isSelfMint ? 'bg-purple-600 text-white' : 'bg-gray-200 text-gray-700'}`}
              >
                Other
              </button>
            </div>
          </div>
          {!isSelfMint && (
            <div>
              <label htmlFor={`recipientAddress-${title}`} className="block text-sm font-medium text-gray-700 mb-1">
                Recipient Address
              </label>
              <input
                type="text"
                id={`recipientAddress-${title}`}
                value={recipientAddress}
                onChange={(e) => setRecipientAddress(e.target.value)}
                className="block text-purple-300 w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50 transition duration-200"
                placeholder="Enter recipient address"
              />
            </div>
          )}
          <motion.button
            type="button"
            className="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-opacity-50 flex items-center justify-center"
            onClick={handleAction}
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
          >
            {functionName}
            <motion.div
              className="ml-2"
              initial={{ x: -5, opacity: 0 }}
              animate={{ x: isHovered ? 0 : -5, opacity: isHovered ? 1 : 0 }}
              transition={{ duration: 0.2 }}
            >
              <ArrowRight className="w-4 h-4" />
            </motion.div>
          </motion.button>
        </form>
      </div>
      <motion.div 
        className="bg-purple-100 px-6 py-4 text-sm text-purple-600"
        initial={{ opacity: 0, height: 0 }}
        animate={{ opacity: isHovered ? 1 : 0, height: isHovered ? 'auto' : 0 }}
        transition={{ duration: 0.3 }}
      >
        Contract: {contractAddress.slice(0, 6)}...{contractAddress.slice(-4)}
      </motion.div>
    </motion.div>
  );
}