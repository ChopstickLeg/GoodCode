import { Helmet } from "react-helmet";
import { useParams } from "react-router-dom";

const PrPage = () => {
  const { prID } = useParams<{ prID: string }>();
  return (
    <div>
      <Helmet>
        <title>PR Page - {prID}</title>
      </Helmet>
      <div className="flex flex-col items-center justify-center h-screen">
        <h1 className="text-2xl font-bold mb-4">PR Page</h1>
        <p className="text-lg">Pull Request ID: {prID}</p>
        <PRProgress currentStep={parseInt(prID ?? "0")} />
      </div>
    </div>
  );
};
const PRProgress = ({ currentStep }: { currentStep: number }) => {
  const steps = [
    "topic -> feature",
    "code in feature",
    "feature -> master-next",
    "code in master-next",
  ];
  const width = "w-[800px]";
  const nodeSize = "w-6 h-6";
  const lineHeight = "h-1";
  const spacing = "w-52";

  return (
    <div className="mt-4">
      <div className={`relative pt-1 ${width}`}>
        <div className="flex mb-2 items-center justify-between">
          {steps.map((step, index) => (
            <div key={index} className="relative">
              <div
                className={`${nodeSize} rounded-full ${
                  currentStep === 3
                    ? "bg-green-500"
                    : index <= currentStep
                    ? "bg-yellow-500"
                    : "bg-gray-200"
                }`}
              />
              <div className="text-sm mt-2 whitespace-nowrap">{step}</div>
              {index < steps.length - 1 && (
                <div
                  className={`absolute top-[11px] -right-52 ${lineHeight} ${spacing} ${
                    currentStep === 3
                      ? "bg-green-500"
                      : index < currentStep
                      ? "bg-yellow-500"
                      : "bg-gray-200"
                  }`}
                />
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
export default PrPage;
