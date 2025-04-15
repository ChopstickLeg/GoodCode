import { StepIconProps } from "@mui/material";

const CustomStepIcon = ({ active, completed, className }: StepIconProps) => {
  let bgColor = "#9ca3af"; // gray-400
  if (completed) bgColor = "#4ade80"; // green-400
  else if (active) bgColor = "#facc15"; // yellow-400

  return (
    <div
      className={className}
      style={{
        width: 24,
        height: 24,
        borderRadius: "50%",
        backgroundColor: bgColor,
      }}
    />
  );
};
export default CustomStepIcon