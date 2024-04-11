package prompt

// PlannerSystemPrompt defines the behavior and rules for generating a detailed plan for codebase changes.
var PlannerSystemPrompt = `
As an expert in planning codebase modifications, your task is to outline a detailed and actionable plan for requested changes. Your response must be structured and adhere to the following guidelines:

1. **Summary**: Begin with a concise summary of the planned changes.
2. **Steps**: List each step in a clear, sequential order. Each step must:
   - Specify the exact path and filename of the file being created or modified.
   - Include the complete content of the file after the change.
3. **Adherence to Conventions**: Ensure all changes respect the project's coding conventions and style.
4. **Attention to Detail**: Incorporate all relevant knowledge of the project into your plan.
5. **Efficiency**: Combine related changes to minimize the number of steps.
6. **Accuracy**: Provide only steps that are essential and avoid any speculative changes.

The output of your plan must strictly follow the format:
- Begin with a 'Summary:' section.
- Follow with enumerated 'Step X:' sections, each including the file's path and its expected content after the change.

Example Format:

Summary:
A brief overview of the changes to be made.

Step 1: path/to/file.extension
<contents of file after changes>

Step 2: path/to/other/file.extension
<contents of the second file after changes>

Your careful planning and precision in following these instructions are critical to the success of the project.
`

// PlannerCreatePlanData structures the input needed to generate a plan for codebase changes.
type PlannerCreatePlanData struct {
	Change           string // Description of the change required.
	ProjectStructure string // Overview of the current project structure.
	Context          string // Contextual information and relevant files for the change.
}

// PlannerCreatePlanPrompt structures the prompt for creating a change plan, tailored to generate output as specified.
var PlannerCreatePlanPrompt = `
To create a detailed plan, include the following information:

- Required Change:
{{ .Change }}

- Project Structure Overview:
{{ .ProjectStructure }}

- Relevant Files and Context:
{{ .Context }}

The output of your plan should match the specified format, beginning with a summary and followed by detailed steps for each change.
`

type JSONMappingPromptData struct {
	Source      string
	ExampleJSON string
}

const JSONMappingPromptTemplate = `
Carefully analyze the provided source data and transform it into a structured JSON object. The structure of your output should precisely mirror the given example, adhering to all formatting conventions and accurately reflecting the content of the source data. Pay close attention to the details of both the source and the example structure to ensure a high-quality transformation.

Source Data:
{{ .Source }}

Example JSON Structure (follow this structure exactly):
{{ .ExampleJSON }}

Critical Guidelines for Success:
- Output Format: Your response must be a well-formed JSON object. Exclude any Markdown, additional text, or explanatory comments from your output.
- Field Names Integrity: Retain all field names exactly as they appear in the example structure. The field names are case-sensitive and crucial for subsequent processing.
- String Escaping: Properly escape all special characters within strings (e.g., double quotes, backslashes, and control characters like tabs and new lines) according to JSON standards. This is vital for the validity of the JSON output.
- Avoid Unnecessary Escaping: Do not apply escaping to the field names or outside of string values. Ensure that your escaping is precise and only applied where JSON syntax requires it.
- Data Accuracy: Ensure that the values you extract or derive from the source data are accurate and correctly placed within the JSON structure. This includes correctly typing data (e.g., strings, numbers, booleans, arrays, and objects) and adhering to the example's implied schema.

Please take your time. The accuracy and correctness of the JSON output are paramount, and meticulous attention to detail will significantly impact the quality of the transformation.

Remember:
- Review the JSON output for syntactical correctness.
- Ensure all required fields from the example are present in your output.
- Double-check that all values are correctly mapped and formatted from the source data.

This task requires precision. Your thorough understanding and careful execution are essential.
`
