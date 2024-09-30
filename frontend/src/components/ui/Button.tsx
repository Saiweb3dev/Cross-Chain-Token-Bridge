import React from 'react'
import Link from 'next/link'

interface ButtonProps {
  text: string
  href?: string
  onClick?: () => void
  className?: string
}

const Button: React.FC<ButtonProps> = ({ text, href, onClick, className = '' }) => {
  const baseClasses = `
    inline-flex items-center justify-center px-4 py-2 rounded-md
    text-white bg-purple-800 border-2 border-purple-800
    hover:text-white hover:bg-purple-600 hover:border-white
    transition-colors duration-200 ease-in-out
    focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500
    ${className}
  `

  if (href) {
    return (
      <Link href={href} className={baseClasses}>
        {text}
      </Link>
    )
  }

  return (
    <button
      className={baseClasses}
      onClick={onClick}
      type={onClick ? 'button' : 'submit'}
    >
      {text}
    </button>
  )
}

export default Button