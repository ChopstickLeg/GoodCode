import { Stepper, Step, StepLabel } from "@mui/material";
import CustomStepIcon from "./CustomStepIcon";
import ColorConnector from "./CustomStepConnecter";

const PRProgress = ({ currentStep }: { currentStep: number }) => {
  const steps = [
    "topic ➝ feature",
    "code in feature",
    "feature ➝ master-next",
    "code in master-next",
  ];

  return (
    <div className="bg-[#2a2a2a] rounded-xl px-8 py-4 shadow-md w-full">
      <Stepper
        activeStep={currentStep}
        connector={<ColorConnector />}
        alternativeLabel={false}
      >
        {steps.map((label, index) => {
          // Determine the color based on the step's status
          let labelColor = "#ffffff"; // Default to white for upcoming steps
          if (index < currentStep) {
            labelColor = "#22c55e"; // Green for completed steps
          } else if (index === currentStep) {
            labelColor = "#facc15"; // Yellow for the current step
          }

          return (
            <Step key={index}>
              <StepLabel
                StepIconComponent={CustomStepIcon}
                sx={{
                  "& .MuiStepLabel-label": {
                    color: labelColor,
                    fontSize: "0.875rem",
                    fontFamily: "monospace",
                  },
                }}
              >
                {label}
              </StepLabel>
            </Step>
          );
        })}
      </Stepper>
    </div>
  );
};

export default PRProgress;
