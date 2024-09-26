"use client"
import React, { useState, useEffect } from 'react'
import Link from 'next/link'
import { notFound } from 'next/navigation'

interface ContractData {
  name: string;
  description: string;
  details: string;
  contractAddress: string;
  abi: any[];
}

export default function ContractPage({ params }: { params: { index: string } }) {
  const [contract, setContract] = useState<ContractData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchContractData = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`/api/contracts/${params.index}`);
        if (!response.ok) {
          throw new Error('Failed to fetch contract data');
        }
        const data = await response.json();
        setContract(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setIsLoading(false);
      }
    };

    fetchContractData();
  }, [params.index]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error || !contract) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-md overflow-hidden">
        <div className="p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">
            {contract.name} Contract Details
          </h1>
          <p className="text-gray-600 mb-4">
            {contract.description}
          </p>
          <p className="text-gray-600 mb-4">
            {contract.details}
          </p>
          <p className="text-gray-600 mb-4">
            Contract Address: {contract.contractAddress}
          </p>
          <details>
            <summary className="text-blue-500 hover:text-blue-600 cursor-pointer">View ABI</summary>
            <pre className="mt-2 p-4 bg-gray-100 rounded overflow-x-auto">
              {JSON.stringify(contract.abi, null, 2)}
            </pre>
          </details>
          <Link href="/Contracts" className="mt-4 inline-block text-blue-500 hover:text-blue-600 font-medium">
            ← Back to all contracts
          </Link>
        </div>
      </div>
    </div>
  )
}