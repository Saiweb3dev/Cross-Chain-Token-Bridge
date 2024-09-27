import React from 'react'
import TokenCard from '@/components/ui/Contract_UI/ContractCard'

const contracts = [
  {
    title: "Token",
    caption: "Bitcoin futures contract allows traders to speculate on the future price of Bitcoin without owning the underlying asset.",
    href: "/Contracts/Token"
  },
  {
    title: "Vault",
    caption: "Ethereum futures contract enables traders to bet on the future price of Ether, the native cryptocurrency of the Ethereum network.",
    href: "/Contracts/Vault"
  },
  {
    title: "Router",
    caption: "Tether futures contract provides a way to trade on the stability of USDT against other cryptocurrencies or fiat currencies.",
    href: "/Contracts/Router"
  }
]

export default function ContractsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-100 to-purple-200 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-5xl font-extrabold text-purple-800 text-center mb-10">
          Select Contract to interact
        </h1>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {contracts.map((contract, index) => (
            <TokenCard
              key={index}
              title={contract.title}
              caption={contract.caption}
              href={contract.href}
            />
          ))}
        </div>
      </div>
    </div>
  )
}