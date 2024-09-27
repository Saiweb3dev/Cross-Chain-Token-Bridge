"use client"

import React, { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Check, ChevronDown } from 'lucide-react'
import { useChain } from '../../contexts/chainContext'

type ChainOption = {
  id: string
  name: string
  caption: string
}

const chainOptions: ChainOption[] = [
  { id: '11155111', name: 'Sepolia', caption: 'ID: 11155111' },
  { id: '80002', name: 'Amoy', caption: 'ID: 80002' },
  { id: '5', name: 'Goerli', caption: 'ID: 5' },
]

export default function ChainSwitcher() {
  const chainContext = useChain()
  const [isOpen, setIsOpen] = useState(false)

  if (!chainContext) {
    return <div className="animate-pulse bg-gray-200 h-10 w-32 rounded-md"></div>
  }

  const { chainId, setChainId } = chainContext

  const handleChainSwitch = (newChainId: string) => {
    setChainId(newChainId)
    setIsOpen(false)
  }

  const currentChain = chainOptions.find(option => option.id === chainId) || chainOptions[0]

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="bg-white border border-gray-200 text-gray-800 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-50 transition duration-300 flex items-center space-x-2 shadow-sm"
      >
        <span>{currentChain.name}</span>
        <ChevronDown className={`w-4 h-4 transition-transform duration-300 ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
            onClick={() => setIsOpen(false)}
          >
            <motion.div
              className="bg-white p-6 rounded-lg shadow-xl max-w-sm w-full m-4"
              onClick={(e) => e.stopPropagation()}
            >
              <h2 className="text-2xl font-bold mb-4">Select Chain</h2>
              <div className="grid gap-4">
                {chainOptions.map((option) => (
                  <motion.button
                    key={option.id}
                    onClick={() => handleChainSwitch(option.id)}
                    className={`p-4 rounded-lg text-left transition duration-300 flex items-center justify-between ${
                      chainId === option.id
                        ? 'bg-purple-100 text-purple-700'
                        : 'bg-gray-100 text-gray-800 hover:bg-gray-200'
                    }`}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                  >
                    <div>
                      <div className="font-semibold">{option.name}</div>
                      <div className="text-sm text-gray-500">{option.caption}</div>
                    </div>
                    {chainId === option.id && (
                      <Check className="w-5 h-5 text-purple-700" />
                    )}
                  </motion.button>
                ))}
              </div>
              <button
                onClick={() => setIsOpen(false)}
                className="mt-4 bg-gray-300 text-gray-800 px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-400 transition duration-300 w-full"
              >
                Close
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}