// adapted from https://observablehq.com/@tezzutezzu/world-history-timeline
"use strict"


function createTooltip(el) {
  el
    .style("position", "absolute")
    .style("pointer-events", "none")
    .style("top", 0)
    .style("opacity", 0)
    .style("background", "white")
    .style("border-radius", "5px")
    .style("box-shadow", "0 0 10px rgba(0,0,0,.25)")
    .style("padding", "10px")
    .style("line-height", "1.3")
    .style("font-size", "0.7rem")
    .style("color", "black")
}

function getTooltipContent(d) {
  return `<strong>${d.Title}</strong>
<br/>
<b style="color:${d.color.darker()}">${d.Organizer}</b>
<br/>
${d.StartDate} - ${d.EndDate}
`
}

function drawChart(data) {
  let projects = data.Projects
  const parent = d3.select("#chart");

  const width = parent.node().clientWidth
  const height = 50 + 30*projects.length
  const margin = ({top: 30, right: 30, bottom: 30, left: 30})

  const x = d3.scaleTime()
    .domain([d3.min(projects, d => d.StartDate), d3.max(projects, d => d.EndDate)].map(d => new Date(d)))
    .range([0, width - margin.left - margin.right])

  const y = d3.scaleBand()
      .domain(d3.range(projects.length))
      .range([0,height - margin.bottom - margin.top])
      .padding(0.2)

  const axisBottom = d3.axisBottom(x)
    .tickPadding(2)

  const axisTop = d3.axisTop(x)
    .tickPadding(2)

  const color = d3.scaleOrdinal(d3.schemeSet3).domain(projects.map(p => p.Organizer))

  function getRect(d) {
    const el = d3.select(this);
    const sx = x(new Date(d.StartDate));
    const w = x(new Date(d.EndDate)) - x(new Date(d.StartDate));
    // Put label on the side with more space.
    const isLabelLeft = sx > width - (sx+w)
    console.log(d.Title, sx, w, isLabelLeft)

    el.style("cursor", "pointer")

    el
      .append("rect")
      .attr("x", sx)
      .attr("height", y.bandwidth())
      .attr("width", w)
      .attr("fill", d.color);

    el
      .append("text")
      .text(d.Title)
      .attr("x", isLabelLeft ? sx-5 : sx+w+5)
      .attr("y", 2.5)
      .attr("fill", "currentColor")
      .style("text-anchor", isLabelLeft ? "end" : "start")
      .style("dominant-baseline", "hanging");
  }

  projects.forEach(d => d.color = d3.color(color(d.Organizer)))

  const svg = parent.append("svg")
    .attr("width", width)
    .attr("height", height)
    .attr("viewBox", `0,0,${width},${height}`)

  const g = svg.append("g").attr("transform", (d,i)=>`translate(${margin.left} ${margin.top})`);

  const groups = g
    .selectAll("g")
    .data(projects)
    .enter()
    .append("g")
    .attr("class", "project")


  const tooltip = d3.select(document.createElement("div")).call(createTooltip);

  const line = svg.append("line").attr("y1", margin.top-10).attr("y2", height-margin.bottom).attr("stroke", "rgba(129,210,199,0.7)").style("pointer-events","none");

  groups.attr("transform", (d,i)=>`translate(0 ${y(i)})`)

  groups
    .each(getRect)
    .on("mouseover", function(event, d) {
      d3.select(this).select("rect").attr("fill", d.color.darker())

      tooltip
        .style("opacity", 1)
        .html(getTooltipContent(d))
    })
    .on("mouseleave", function(event, d) {
      d3.select(this).select("rect").attr("fill", d.color)
      tooltip.style("opacity", 0)
    })


  svg
    .append("g")
    .attr("transform", (d,i)=>`translate(${margin.left} ${margin.top-10})`)
    .call(axisTop)

  svg
    .append("g")
    .attr("transform", (d,i)=>`translate(${margin.left} ${height-margin.bottom})`)
    .call(axisBottom)



  svg.on("mousemove", function(event) {

    let [x,y] = d3.pointer(event);
    line.attr("transform", `translate(${x} 0)`);
    y +=20;
    if(x>width/2) x-= 100;

    tooltip
      .style("left", x + "px")
      .style("top", y + "px")
  })

  parent.node().appendChild(svg.node());
  parent.node().appendChild(tooltip.node());

}

d3.json("projects.json").then(drawChart)
