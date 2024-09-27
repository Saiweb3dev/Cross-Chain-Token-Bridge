import React from 'react'
import { motion } from 'framer-motion'
import Link from 'next/link'
import { AlertCircle } from 'lucide-react'

interface ContractErrorProps {
  error: string;
}

export default function ContractError({ error }: ContractErrorProps) {
  return (
    <motion.div 
      initial={{ opacity: 0 }} 
      animate={{ opacity: 1 }} 
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-red-50 to-red-100 py-12 px-4 sm:px-6 lg:px-8 flex items-center justify-center"
    >
      <div className="max-w-md w-full bg-white rounded-xl shadow-lg overflow-hidden">
        <div className="p-8">
          <div className="flex items-center justify-center mb-4">
            <AlertCircle className="text-red-500 w-12 h-12" />
          </div>
          <h1 className="text-3xl font-bold text-red-600 mb-4 text-center">Error</h1>
          <p className="text-gray-600 mb-6 text-center">{error}</p>
          <Link href="/Contracts" className="block text-center text-purple-600 hover:text-purple-800 font-medium">
            ← Back to all contracts
          </Link>
        </div>
      </div>
    </motion.div>
  )
}