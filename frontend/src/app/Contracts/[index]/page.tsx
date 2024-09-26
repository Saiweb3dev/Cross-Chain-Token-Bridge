import React from 'react'
import Link from 'next/link'
import { notFound } from 'next/navigation'

const contractData = {
  Custom_Token: {
    name: "Custom Token",
    description: "Bitcoin futures contract allows traders to speculate on the future price of Bitcoin without owning the underlying asset.",
    details: "Bitcoin is the world's first cryptocurrency and remains the largest by market capitalization. Futures contracts for Bitcoin allow investors to gain exposure to BTC price movements without holding the actual cryptocurrency."
  },
  Vault: {
    name: "Vault",
    description: "Ethereum futures contract enables traders to bet on the future price of Ether, the native cryptocurrency of the Ethereum network.",
    details: "Ethereum is a decentralized, open-source blockchain featuring smart contract functionality. ETH futures provide a way for traders to speculate on Ether's price movements or hedge their existing Ethereum holdings."
  },
  Router: {
    name: "Router",
    description: "Tether futures contract provides a way to trade on the stability of USDT against other cryptocurrencies or fiat currencies.",
    details: "Tether (USDT) is a stablecoin pegged to the US dollar. USDT futures contracts allow traders to speculate on the stability of USDT or its relationship with other cryptocurrencies."
  }
}

export default function ContractPage({ params }: { params: { index: string } }) {
  const contract = contractData[params.index as keyof typeof contractData]

  if (!contract) {
    notFound()
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
          <Link href="/Contracts" className="text-blue-500 hover:text-blue-600 font-medium">
            ← Back to all contracts
          </Link>
        </div>
      </div>
    </div>
  )
}