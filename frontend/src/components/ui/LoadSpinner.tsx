import React from 'react'
import { motion } from 'framer-motion'

export default function LoadingSpinner() {
  return (
    <motion.div 
      initial={{ opacity: 0 }} 
      animate={{ opacity: 1 }} 
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-purple-100 to-purple-200 flex justify-center items-center"
    >
      <div className="w-16 h-16 border-4 border-purple-600 border-t-transparent rounded-full animate-spin"></div>
    </motion.div>
  )
}