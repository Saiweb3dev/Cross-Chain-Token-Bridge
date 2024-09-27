"use client"

import React, { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import ContractDetails from '@/components/ui/Contract_UI/ContractDetails'
import ContractFunction from '@/components/ui/Contract_UI/ContractFunction'
import ContractError from '@/components/ui/Contract_UI/ContractError'
import LoadingSpinner from '@/components/ui/LoadSpinner'
import { useChain } from '@/contexts/chainContext'
interface ContractData {
  name: string;
  description: string;
  details: string;
  contractAddress: string;
  abi: any[];
}
 // Fixed chainID

async function getContractData(index: string, chainId: string): Promise<ContractData> {
 
  const res = await fetch(`http://localhost:8080/api/contract/${chainId}/${index}`, { cache: 'no-store' });
  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.error || 'Failed to fetch contract data');
  }
  return res.json();
}

export default function ContractPage({ params }: { params: { index: string } }) {
  const { chainId } = useChain();
  const [contract, setContract] = useState<ContractData | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    setError(null);
    getContractData(params.index, chainId)
      .then(data => {
        setContract(data);
        setLoading(false);
      })
      .catch(err => {
        setError(err instanceof Error ? err.message : 'An error occurred');
        setLoading(false);
      });
  }, [params.index, chainId]);

  if (error) {
    return <ContractError error={error} />;
  }

  if (!contract) {
    return <LoadingSpinner />;
  }

  return (
    <motion.div 
      initial={{ opacity: 0 }} 
      animate={{ opacity: 1 }} 
      exit={{ opacity: 0 }}
      className="min-h-screen bg-gradient-to-br from-purple-100 to-purple-200 py-12 px-4 sm:px-6 lg:px-8"
    >
      <div className="max-w-4xl mx-auto space-y-8">
        <ContractDetails contract={contract} chainId={chainId} />
        <motion.div 
      className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ staggerChildren: 0.1 }}
    >
      <ContractFunction 
        title={`Mint ${contract.name} token`}
        contractAddress={contract.contractAddress}
        abi={contract.abi}
      />
      <ContractFunction 
        title={`Burn ${contract.name} token`}
        contractAddress={contract.contractAddress}
        abi={contract.abi}
      />
    </motion.div>
      </div>
    </motion.div>
  )
}