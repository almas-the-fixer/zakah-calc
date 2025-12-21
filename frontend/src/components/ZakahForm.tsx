import { useState } from 'react';
import axios from 'axios';

// Define the shape of the response
interface CalculationResponse {
  total_assets: number;
  nisab_threshold: number;
  zakah_payable: number;
  currency: string;
  message: string;
}

export default function ZakahForm() {
  const [formData, setFormData] = useState({
    gold_grams: 0,
    silver_grams: 0,
    cash: 0,
    business_assets: 0,
    liabilities: 0,
    currency: 'USD', // Default
  });

  const [result, setResult] = useState<CalculationResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: name === 'currency' ? value : parseFloat(value) || 0,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setResult(null);

    try {
      // 🚀 Call backend
      const response = await axios.post('http://localhost:8080/calculate-zakah', formData);
      setResult(response.data);
    } catch (err) {
      console.error(err);
      setError('Failed to connect to backend. Is it running on port 8080?');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-xl mx-auto mt-10 p-8 bg-white rounded-2xl shadow-xl border border-gray-100">
      <h2 className="text-3xl font-bold text-center text-emerald-600 mb-8">Zakah Calculator</h2>

      <form onSubmit={handleSubmit} className="space-y-5">
        {/* Currency Selector */}
        <div>
          <label className="block text-sm font-semibold text-gray-700 mb-1">Currency</label>
          <select
            name="currency"
            value={formData.currency}
            onChange={handleChange}
            className="w-full p-3 bg-gray-50 border border-gray-200 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 outline-none transition"
          >
            <option value="USD">USD ($)</option>
            <option value="INR">INR (₹)</option>
            <option value="PKR">PKR (Rs)</option>
            <option value="EUR">EUR (€)</option>
            <option value="GBP">GBP (£)</option>
          </select>
        </div>

        {/* Inputs Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          <InputGroup label="Gold (grams)" name="gold_grams" onChange={handleChange} />
          <InputGroup label="Silver (grams)" name="silver_grams" onChange={handleChange} />
          <InputGroup label="Cash Assets" name="cash" onChange={handleChange} />
          <InputGroup label="Business Assets" name="business_assets" onChange={handleChange} />
          <InputGroup label="Liabilities (Debt)" name="liabilities" onChange={handleChange} />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-emerald-600 text-white py-3.5 rounded-lg font-bold text-lg hover:bg-emerald-700 active:scale-95 transition-all disabled:opacity-50 disabled:scale-100"
        >
          {loading ? 'Calculating...' : 'Calculate Zakah'}
        </button>
      </form>

      {/* Error Message */}
      {error && (
        <div className="mt-6 p-4 bg-red-50 text-red-600 rounded-lg text-center border border-red-100">
          {error}
        </div>
      )}

      {/* Results Section */}
      {result && (
        <div className="mt-8 p-6 bg-emerald-50 rounded-xl border border-emerald-100 animation-fade-in">
          <div className="text-center mb-6">
            <h3 className="text-xl font-bold text-emerald-800">{result.message}</h3>
          </div>
          
          <div className="grid grid-cols-2 gap-4 mb-6">
            <ResultCard label="Net Assets" value={result.total_assets} currency={result.currency} />
            <ResultCard label="Nisab Threshold" value={result.nisab_threshold} currency={result.currency} />
          </div>

          <div className="bg-emerald-600 text-white p-5 rounded-lg text-center shadow-lg transform transition hover:scale-105">
            <p className="text-sm font-medium opacity-90 uppercase tracking-wider">Zakah Payable</p>
            <p className="text-4xl font-extrabold mt-1">
              {result.zakah_payable.toLocaleString(undefined, { maximumFractionDigits: 2 })} {result.currency}
            </p>
          </div>
        </div>
      )}
    </div>
  );
}

// Helper Components for clean code
const InputGroup = ({ label, name, onChange }: any) => (
  <div>
    <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
    <input
      type="number"
      name={name}
      onChange={onChange}
      placeholder="0"
      className="w-full p-3 bg-gray-50 border border-gray-200 rounded-lg focus:ring-2 focus:ring-emerald-500 outline-none transition"
    />
  </div>
);

const ResultCard = ({ label, value, currency }: any) => (
  <div className="bg-white p-4 rounded-lg shadow-sm border border-emerald-100 text-center">
    <p className="text-xs font-semibold text-gray-400 uppercase">{label}</p>
    <p className="text-lg font-bold text-gray-800 mt-1">
      {value.toLocaleString(undefined, { maximumFractionDigits: 0 })} {currency}
    </p>
  </div>
);