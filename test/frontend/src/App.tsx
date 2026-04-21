import { useEffect, useState } from 'react'
import type { OilPrice } from './types/oil'
import { fetchPrices, pullPrices } from './services/api'
import './App.css'

function App() {
  const [prices, setPrices] = useState<OilPrice[]>([])
  const [loading, setLoading] = useState(false)
  const [filter, setFilter] = useState('')
  const [error, setError] = useState<string | null>(null)

  const loadData = async (type?: string) => {
    setLoading(true)
    setError(null)
    try {
      const data = await fetchPrices(type)
      setPrices(data || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  const handlePull = async () => {
    setLoading(true)
    setError(null)
    try {
      await pullPrices()
      await loadData(filter)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData(filter)
  }, [filter])

  return (
    <div className="container">
      <header>
        <h1>Bangchak Oil Prices</h1>
      </header>
      
      <main>
        <div className="controls">
          <div className="search-box">
            <input 
              type="text" 
              placeholder="Filter by oil type (e.g. Gasohol)..." 
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
          <button 
            className="btn-primary" 
            onClick={handlePull} 
            disabled={loading}
          >
            {loading ? 'Updating...' : 'Pull Latest Prices'}
          </button>
        </div>

        {error && <div className="error-message">Error: {error}</div>}
        
        <div className="table-container">
          {loading && !prices.length ? (
            <div className="loading">Loading prices...</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Oil Type</th>
                  <th>Price (THB)</th>
                  <th>Effective Date</th>
                </tr>
              </thead>
              <tbody>
                {prices.length > 0 ? (
                  prices.map((p) => (
                    <tr key={p.id}>
                      <td className="oil-type">{p.oil_type}</td>
                      <td className="price">{p.price.toFixed(2)}</td>
                      <td className="date">{new Date(p.date).toLocaleDateString()}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={3} className="no-data">No oil prices found</td>
                  </tr>
                )}
              </tbody>
            </table>
          )}
        </div>
      </main>
      
      <footer>
        <p>Data source: Bangchak API</p>
      </footer>
    </div>
  )
}

export default App
