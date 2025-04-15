import { Helmet } from "react-helmet";
import { useParams } from "react-router-dom";
import PRProgress from "./components/PRProgress";

const PrPage = () => {
  const { prID } = useParams<{ prID: string }>();
  const currentStep = parseInt(prID ?? "0");

  return (
    <div className="min-h-screen bg-[#242424] text-white">
      <Helmet>
        <title>PR Page - {prID}</title>
      </Helmet>

      <div className="w-full max-w-7xl mx-auto px-6 py-6">
        <div className="bg-[#333] p-6 rounded-xl shadow-md">
          <PRProgress currentStep={currentStep} />
        </div>
      </div>

      {/* Two-Column Grid: PR Info + Comments */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-7xl mx-auto px-6 py-6">
        <section className="bg-[#333] p-6 rounded-xl shadow-md">
          <h2 className="text-xl font-semibold mb-2">Pull Request Info</h2>
          <p>Pull Request ID: {prID}</p>
          <p>Details about the PR go here. Maybe author, branch info, etc.</p>
        </section>

        <section className="bg-[#333] p-6 rounded-xl shadow-md">
          <h2 className="text-xl font-semibold mb-2">Comments</h2>
          <p>
            Comment section placeholder. Or maybe no one cared enough to comment
            yet.
          </p>
        </section>

        {/* Code Diff + AI Roast */}
        <section className="bg-[#333] p-6 rounded-xl shadow-md md:col-span-1">
          <h2 className="text-xl font-semibold mb-2">Code Diff</h2>
          <p>This will be the diff viewer. Coming up next...</p>
        </section>

        <section className="bg-[#333] p-6 rounded-xl shadow-md md:col-span-1">
          <h2 className="text-xl font-semibold mb-2">AI Roast</h2>
          <p>
            This will be the AI Roasting section. Coming to a browser near
            you...
          </p>
        </section>
      </div>
    </div>
  );
};

export default PrPage;
