import { NextApiRequest, NextApiResponse } from 'next'

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { index } = req.query

  if (req.method !== 'GET') {
    return res.status(405).json({ message: 'Method not allowed' })
  }

  try {
    // Replace with your actual Go backend URL
    const backendUrl = `http://localhost:8080/api/contracts/${index}`
    const response = await fetch(backendUrl)

    if (!response.ok) {
      throw new Error('Failed to fetch contract data from backend')
    }

    const data = await response.json()
    res.status(200).json(data)
  } catch (error) {
    console.error('Error fetching contract data:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
}