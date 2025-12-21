import ZakahForm from './components/ZakahForm';

function App() {
  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 flex flex-col items-center">
      <div className="max-w-4xl w-full space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-extrabold text-gray-900 tracking-tight sm:text-5xl">
            Zakah<span className="text-emerald-600">Calc</span>
          </h1>
          <p className="mt-3 text-lg text-gray-500">
            Secure, Private, and Real-time Zakah Calculation.
          </p>
        </div>
        
        <ZakahForm />
        
        <footer className="mt-16 text-center text-gray-400 text-sm">
          <p>&copy; 2025 ZakahCalc Project. Built with Go & React.</p>
        </footer>
      </div>
    </div>
  );
}

export default App;