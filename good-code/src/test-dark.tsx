// Test component to force Tailwind to include dark mode classes
export const TestDark = () => {
  return (
    <div className="bg-white dark:bg-gray-900 text-black dark:text-white p-4">
      <h1 className="text-blue-500 dark:text-blue-400">Test Dark Mode</h1>
      <p className="bg-gray-100 dark:bg-gray-800">This should test dark mode</p>
    </div>
  );
};
