import React from 'react'
import Link from 'next/link'

interface TokenCardProps {
  title: string
  caption: string
  href: string
}

export default function TokenCard({ title, caption, href }: TokenCardProps) {
  return (
    <Link href={href} className="block w-full" aria-label={`View details for ${title} contract`}>
      <div className="h-64 rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300 ease-in-out bg-white cursor-pointer transform hover:scale-105">
        <div className="h-full flex flex-col justify-between p-6">
          <div className="text-center">
            <h2 className="font-bold text-2xl sm:text-3xl md:text-4xl mb-2 text-purple-800">{title}</h2>
          </div>
          <p className="text-gray-700 text-sm sm:text-base text-center">{caption}</p>
        </div>
      </div>
    </Link>
  )
}