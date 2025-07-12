import Header from "../Header";

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-gray-50 to-blue-50 text-gray-900">
      <Header />
      <main className="w-full p-6 max-w-7xl mx-auto">
        <div className="animate-fade-in">{children}</div>
      </main>
    </div>
  );
};

export default Layout;
