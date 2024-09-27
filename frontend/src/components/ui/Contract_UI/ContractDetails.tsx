"use client"

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { ChevronDown, ChevronUp, X, Eye } from 'lucide-react'
import Link from 'next/link'

interface ContractDetailsProps {
  contract: {
    name: string;
    description: string;
    details: string;
    contractAddress: string;
    abi: any[];
  };
  chainId: string;
}

export default function ContractDetails({ contract, chainId }: ContractDetailsProps) {
  const [isFullAbiOpen, setIsFullAbiOpen] = useState(false);

  const fullAbi = JSON.stringify(contract.abi, null, 2);

  return (
    <motion.div 
      className="bg-white rounded-xl shadow-lg overflow-hidden"
      initial={{ y: 20 }}
      animate={{ y: 0 }}
      transition={{ type: "spring", stiffness: 260, damping: 20 }}
    >
      <div className="p-8">
        <h1 className="text-4xl font-bold text-purple-800 mb-4">
          {contract.name} Contract Details
        </h1>
        <p className="text-gray-600 mb-4 text-lg">
          {contract.description}
        </p>
        <p className="text-gray-600 mb-6">
          {contract.details}
        </p>
        <div className="bg-purple-100 rounded-lg p-4 mb-6">
          <p className="text-purple-800 mb-2">
            <span className="font-semibold">Contract Address:</span> {contract.contractAddress}
          </p>
          <p className="text-purple-800">
            <span className="font-semibold">Chain ID:</span> {chainId}
          </p>
        </div>
        
        <div className="space-y-4">
          <button 
            onClick={() => setIsFullAbiOpen(true)}
            className="flex items-center text-purple-600 hover:text-purple-800 font-medium focus:outline-none"
          >
            <Eye className="mr-2" />
            View Full ABI
          </button>
        </div>

        <AnimatePresence>
          {isFullAbiOpen && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
            >
              <motion.div
                initial={{ scale: 0.9 }}
                animate={{ scale: 1 }}
                exit={{ scale: 0.9 }}
                className="bg-white rounded-lg p-6 w-full max-w-3xl max-h-[80vh] overflow-y-auto relative"
              >
                <button 
                  onClick={() => setIsFullAbiOpen(false)}
                  className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
                >
                  <X size={24} />
                </button>
                <h2 className="text-2xl font-bold text-purple-800 mb-4">Full ABI</h2>
                <pre className="p-4 bg-purple-800 text-white rounded-lg overflow-x-auto text-sm">
                  {fullAbi}
                </pre>
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>

        <Link href="/Contracts" className="mt-6 inline-flex items-center text-purple-600 hover:text-purple-800 font-medium">
          ← Back to all contracts
        </Link>
      </div>
    </motion.div>
  )
}