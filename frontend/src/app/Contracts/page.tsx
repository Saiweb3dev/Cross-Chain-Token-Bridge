"use client";
import React from 'react'
import TokenCard from '@/components/ui/Contract_UI/ContractCard'
import { useChain } from '@/contexts/chainContext';
import Button from '@/components/ui/Button';

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
  const { chainId, isNetworkSwitchRequired, switchNetwork } = useChain();
  if (isNetworkSwitchRequired) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-100">
        <div className="text-center p-8 bg-white rounded-lg shadow-xl max-w-md w-full">
          <h2 className="text-3xl font-bold text-red-600 mb-6">Network Mismatch</h2>
          <p className="text-gray-700 mb-8">Please switch to the correct network to interact with this contract.</p>
          <Button 
            text="Switch Network"
            onClick={() => switchNetwork(chainId)}
            className="w-full hover:bg-purple-900 text-white"
          />
        </div>
      </div>
    );
  }
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