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
          let labelColor = "#ffffff"; 
          if (index < currentStep) {
            labelColor = "#22c55e"; 
          } else if (index === currentStep) {
            labelColor = "#facc15"; 
          }

          return (
            <Step key={index}>
              <StepLabel
                StepIconComponent={CustomStepIcon}
                sx={{
                  "& .MuiStepLabel-label": {
                    color: `${labelColor} !important`,
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
