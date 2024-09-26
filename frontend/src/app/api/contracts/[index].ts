import { NextApiRequest, NextApiResponse } from 'next'

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { index } = req.query

  if (req.method !== 'GET') {
    return res.status(405).json({ message: 'Method not allowed' })
  }

  try {
    const backendUrl = `http://localhost:8080/api/contracts/${index}`
    console.log('Fetching from:', backendUrl) // Log the URL being fetched

    const response = await fetch(backendUrl)

    if (!response.ok) {
      const errorData = await response.json()
      console.error('Backend error:', errorData) // Log the error from the backend
      throw new Error(errorData.error || 'Failed to fetch contract data from backend')
    }

    const data = await response.json()
    res.status(200).json(data)
  } catch (error) {
    console.error('Error fetching contract data:', error)
    res.status(500).json({ message: error instanceof Error ? error.message : 'Internal server error' })
  }
}