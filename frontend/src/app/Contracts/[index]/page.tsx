import React from 'react'
import Link from 'next/link'

interface ContractData {
  name: string;
  description: string;
  details: string;
  contractAddress: string;
  abi: any[];
}

const CHAIN_ID = '80002'; // Fixed chainID

async function getContractData(index: string): Promise<ContractData> {
  const res = await fetch(`http://localhost:8080/api/contract/${CHAIN_ID}/${index}`, { cache: 'no-store' });
  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.error || 'Failed to fetch contract data');
  }
  return res.json();
}

export default async function ContractPage({ params }: { params: { index: string } }) {
  let contract: ContractData | null = null;
  let error: string | null = null;

  try {
    contract = await getContractData(params.index);
  } catch (err) {
    error = err instanceof Error ? err.message : 'An error occurred';
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto bg-white rounded-xl shadow-md overflow-hidden">
          <div className="p-8">
            <h1 className="text-3xl font-bold text-red-600 mb-4">Error</h1>
            <p className="text-gray-600 mb-4">{error}</p>
            <Link href="/Contracts" className="text-blue-500 hover:text-blue-600 font-medium">
              ← Back to all contracts
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (!contract) {
    return (
      <div>
        <p>Contract not found</p>
        <Link href="/Contracts">← Back to all contracts</Link>
      </div>
    );
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
          <p className="text-gray-600 mb-4">
            Chain ID: {CHAIN_ID}
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