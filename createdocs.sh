#!/bin/bash
tempfile="READMEraw.md"
markdownfile="${tempfile}"
mergedfile='README.md'
mdmerge ${markdownfile} > ${mergedfile}
pdfname='README.pdf'
latexengine='lualatex'
highlightstyle='--highlight-style=zenburn'
# pygments (the default)
# kate
# monochrome
# espresso
# zenburn
# haddock,
# tango.
documentclass='--variable documentclass=ltjarticle'

includeinheader='--include-in-header=fontoptions.tex'
# verbose='--verbose'

echo "building ${pdfname} and ${mergedfile}"
pandoc -f markdown+pandoc_title_block ${verbose} ${highlightstyle} -N --template=template.latex  ${documentclass} --variable papersize:a4 --variable colorlinks  --variable geometry:margin=1in ${includeinheader}  --variable fontsize="12pt" --variable monofont="SourceCodePro-Regular" --variable fontsize="12pt"  ${mergedfile} --latex-engine=${latexengine} --toc -o ${pdfname}

