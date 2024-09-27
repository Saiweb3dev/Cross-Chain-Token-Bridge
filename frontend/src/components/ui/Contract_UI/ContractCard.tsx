"use client"

import React from 'react'
import Link from 'next/link'
import { motion } from 'framer-motion'

interface ContractCardProps {
  title: string
  caption: string
  href: string
}

export default function ContractCard({ title, caption, href }: ContractCardProps) {
  return (
    <Link href={href} className="block w-full" aria-label={`View details for ${title} contract`}>
      <motion.div 
        className="h-64 rounded-xl overflow-hidden shadow-lg bg-white cursor-pointer relative group"
        whileHover={{ scale: 1.05, rotate: 1, boxShadow: '0 0 20px 0 rgba(255, 255, 255, 0.3)' }}
        whileTap={{ scale: 0.95 }}
        transition={{ type: "spring", stiffness: 400, damping: 17 }}
      >
        <motion.div 
          className="absolute inset-0 bg-purple-100 opacity-0 group-hover:opacity-20 transition-opacity duration-300"
          initial={{ opacity: 0 }}
          whileHover={{ opacity: 0.2 }}
        />
        <div className="h-full flex flex-col justify-between p-6 relative z-10">
          <motion.div 
            className="text-center"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
          >
            <h2 className="font-bold text-2xl sm:text-3xl md:text-4xl mb-2 text-purple-800 group-hover:text-purple-900 transition-colors duration-300">
              {title}
            </h2>
          </motion.div>
          <motion.p 
            className="text-purple-600 text-sm sm:text-base text-center group-hover:text-purple-700 transition-colors duration-300"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.4 }}
          >
            {caption}
          </motion.p>
        </div>
        <motion.div 
          className="absolute inset-0 border-2 border-transparent group-hover:border-purple-200 rounded-xl transition-all duration-300"
          initial={{ opacity: 0, scale: 1.1 }}
          whileHover={{ opacity: 1, scale: 1 }}
        />
      </motion.div>
    </Link>
  )
}