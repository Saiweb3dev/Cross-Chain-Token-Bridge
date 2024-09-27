"use client"

import React, { useState } from 'react'
import { motion } from 'framer-motion'
import { ArrowRight } from 'lucide-react'

interface ContractFunctionProps {
  title: string;
  contractAddress: string;
  abi: any[];
}

export default function ContractFunction({ title, contractAddress, abi }: ContractFunctionProps) {
  const [tokenAmount, setTokenAmount] = useState('');
  const [isHovered, setIsHovered] = useState(false);

  const handleAction = () => {
    // Implement contract interaction logic here
    console.log(`Interacting with ${title} for ${tokenAmount} tokens on contract ${contractAddress}`);
  };

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
              className="block w-full px-3 py-2 rounded-md border border-gray-300 shadow-sm focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50 transition duration-200"
              placeholder="Enter amount"
            />
          </div>
          <motion.button
            type="button"
            className="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 transition duration-300 ease-in-out focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-opacity-50 flex items-center justify-center"
            onClick={handleAction}
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
          >
            <span>{title.split(' ')[0]}</span>
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
  )
}