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
`+(d.ReleasedVideo != null ? `<br/>
Release: ${d.ReleasedVideo.Date.replace(/T.*/, '')}
` : "")
}

// countActiveAtDate returns the number of active projects at the given date.
// date is a ISO 8601 date string.
function countActiveAtDate(projects, date) {
  // TODO: Might be more efficient to sort projects?
  let count = 0
  for (let p of projects) {
    if (p.StartDate <= date && date <= p.EndDate)
      count++
  }
  return count
}

function drawChart(projects) {
  const parent = d3.select("#timeline");

  const margin = ({top: 30, right: 30, bottom: 30, left: 30})
  const width = parent.node().clientWidth
  const heightProjects = 50 + 30*projects.length
  const heightCount = 200
  const height = heightProjects + heightCount + 2*margin.bottom

  const x = d3.scaleTime()
    .domain([d3.min(projects, d => d.StartDate), d3.max(projects, d => d.EndDate)].map(d => new Date(d)))
    .range([0, width - margin.left - margin.right])

  const y = d3.scaleBand()
      .domain(d3.range(projects.length))
      .range([0, heightProjects - margin.bottom - margin.top])
      .padding(0.2)

  const axisBottom = d3.axisBottom(x)
    .tickPadding(2)

  const axisTop = d3.axisTop(x)
    .tickPadding(2)

  const color = d3.scaleOrdinal(d3.schemeSet3).domain(projects.map(p => p.Organizer.toLowerCase()))

  function getRect(d) {
    const el = d3.select(this);
    const sx = x(new Date(d.StartDate));
    let w = x(new Date(d.EndDate)) - x(new Date(d.StartDate));

    el.style("cursor", "pointer")

    el
      .append("a")
        .attr("href", d.URL)
        .append("rect")
          .attr("x", sx)
          .attr("height", y.bandwidth())
          .attr("width", w)
          .attr("fill", d.color);

    if (d.ReleasedVideo != null) {
      let xvid = x(new Date(d.ReleasedVideo.Date))
      w += xvid - x(new Date(d.EndDate)) + 10
      el
        .append("a")
          .attr("href", `https://youtu.be/${d.ReleasedVideo.ID}`)
          .append("circle")
            .attr("cx", xvid)
            .attr("cy", y.bandwidth() / 2)
            .attr("r", y.bandwidth() / 5)
            .attr("fill", d.color)
    }

    // Put label on the side with more space.
    const isLabelLeft = sx > width - (sx+w)

    el
      .append("text")
      .text(d.Title.replace(/\([^)]+\)/, ''))
      .attr("x", isLabelLeft ? sx-5 : sx+w+5)
      .attr("y", 2.5)
      .attr("fill", "currentColor")
      .style("text-anchor", isLabelLeft ? "end" : "start")
      .style("dominant-baseline", "hanging");

  }

  // filter change? -> we might need to re-render
  // TODO: make this nicer?
  if (parent.selectAll(".project").size() != projects.length) {
    parent.select("svg").remove()
  }

  // already rendered?
  if (parent.select("svg").empty()) {

    projects.forEach(d => d.color = d3.color(color(d.Organizer.toLowerCase())))

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
        d3.select(this).selectAll("rect, circle").attr("fill", d.color.darker())

        tooltip
          .style("opacity", 1)
          .html(getTooltipContent(d))
      })
      .on("mouseleave", function(event, d) {
        d3.select(this).selectAll("rect, circle").attr("fill", d.color)
        tooltip.style("opacity", 0)
      })


    svg
      .append("g")
      .attr("transform", (d,i)=>`translate(${margin.left} ${margin.top-10})`)
      .call(axisTop)

    svg
      .append("g")
      .attr("transform", (d,i)=>`translate(${margin.left} ${heightProjects-margin.bottom})`)
      .call(axisBottom)


    let yCount = d3.scaleLinear().domain([0, 16]).range([heightCount, 0])
    let axisLeftCount = d3.axisLeft(yCount).tickPadding(2)
    let lineCount = d3.line()
          .x(x)
          .y(d => yCount(countActiveAtDate(projects, toISODateString(d))))
    let gCount = svg
      .append("g")
      .attr("transform", (d,i)=>`translate(${margin.left} ${heightProjects+margin.bottom})`)

    gCount
      .append("path")
        .attr("fill", "none")
        .attr("stroke-width", 1.5)
        .attr("stroke-linejoin", "round")
        .attr("stroke-linecap", "round")
        .attr("stroke", "#81d2c7")
        .attr("d", d => lineCount(x.ticks(100)))

    gCount
      .append("g")
      .call(axisLeftCount)

    gCount
      .append("g")
      .attr("transform", (d,i)=>`translate(0 ${heightCount})`)
      .call(axisBottom)

    gCount
      .append("rect")
        .attr("width", width - margin.right)
        .attr("height", heightCount + margin.bottom)
        .attr("fill", "transparent")
        .on("mousemove", function(event) {
          let [xc,yc] = d3.pointer(event);
          if (xc > margin.left) {
            let date = toISODateString(x.invert(xc))
            let count = countActiveAtDate(projects, date)
            tooltip
              .style("opacity", 1)
              .html(`<strong>${date}</strong>: ${count} active project${count == 1 ? '' : 's'}`)
          }
        })
        .on("mouseleave", function(event) {
            tooltip.style("opacity", 0)
        })


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

  } else {
    // update order
    let prjs = parent.selectAll(".project")

    prjs.data(projects, p => p.Title)
      .transition()
      .ease(d3.easeCubic)
      .attr("transform", (d, i) => `translate(0 ${y(i)})`)

  }

}

// sheetDate converts a date like "October 9th, 2020" to ISO 8601.
function sheetDate(d) {
  const parse = d3.utcParse("%B %d %Y")
  let normalized = d.replace(/([0-9])([a-z]{2})?,\s*/, '$1 ')
                .replace(/\s*\(EXT\)/, '')
                .trim()
  let parsed = parse(normalized)
  if (parsed == null) {
    console.info(`invalid date: ${d} (${normalized})`)
    return null
  }
  return toISODateString(parsed)
}

// dayBefore takes an ISO8601 date and returns the date of the day before.
function dayBefore(d) {
  let date = new Date(d)
  let before = new Date(date.getTime() - 24 * 60 * 60 * 1000)
  return toISODateString(before)
}

// toISODateString converts a Date to an ISO 8601 date string.
function toISODateString(date) {
  return date.toISOString().slice(0, 10)
}

async function main() {
  let [data, sheet] = await Promise.all([d3.json("projects.json"), d3.csv("allprojects.csv")])

  let videosById = new Map(data.Videos.map(v => [v.contentDetails.videoId, v]))

  let sheetProjects = sheet
    .map(p => ({
      FromSheet: true,
      Title: p['Project Name'],
      Organizer: p['Creator'].replace('The Reddit Symphony Orchestra', 'CasuallyNothing').replace(/u\/| .*/g, ''),
      StartDate: sheetDate(p['Start Date']),
      EndDate: sheetDate(p['Deadline']),
      ReleasedVideo: (url => {
        let m = /(youtu\.be\/|youtube\.com\/watch\?v=)([\w-]+)/.exec(url)
        if (m == null) return null
        let v = videosById.get(m[2])
        if (v == null) {
          console.info(`video https://youtu.be/${m[2]} for ${p['Project Name']} not in playlist`)
          return null
        }
        return {
          Title: v.snippet.title,
          ID: v.contentDetails.videoId,
          Date: v.contentDetails.videoPublishedAt,
        }
      })(p['Links to Active Project Page OR Finished Result']),
    }))
    .filter(p => p.StartDate != null && p.EndDate != null && p.ReleasedVideo != null)

  // try to add video links by matching the organizer + start date
  let toidx = p => `${p.Organizer.toLowerCase()} - ${p.StartDate}`
  // hack: due to time zone issues, the start date is sometimes off-by-one
  let toidx2 = p => `${p.Organizer.toLowerCase()} - ${dayBefore(p.StartDate)}`
  let sheetProjectsIndex = new Map(sheetProjects.map(p => [toidx(p), p]))
  data.Projects.forEach(p => {
    let sp
    if ((sp = sheetProjectsIndex.get(toidx(p))) || (sp = sheetProjectsIndex.get(toidx2(p)))) {
      p.ReleasedVideo = sp.ReleasedVideo
    }
  })
  // filter out projects we already know from Reddit posts
  let projectsByVideoId = new Map(data.Projects.map(p => [p.ReleasedVideo?.ID, p]).filter(([k, v]) => k != null))
  sheetProjects = sheetProjects.filter(p => !projectsByVideoId.get(p.ReleasedVideo.ID))

  const allProjects = data.Projects.concat(sheetProjects)
  allProjects.sort((a, b) => d3.ascending(a.EndDate, b.EndDate))

  let projects

  const applySort = () => {
    let by = d3.select("#timeline-sortby").node().value
    let key = by == 'deadline'  ? (p => p.EndDate) :
              by == 'video'     ? (p => p.ReleasedVideo?.Date ?? 'z'+p.EndDate) :
              by == 'organizer' ? (p => p.Organizer.toLowerCase()) :
              console.error("wrong sorby value", by)
    projects.sort((a, b) => d3.ascending(key(a), key(b)))
    drawChart(projects)
  }

  d3.select("#timeline-sortby").on("change", applySort)

  const applyFilter = () => {
    if (d3.select("#timeline-showold").node().checked)
      projects = allProjects
    else
      projects = allProjects.filter(p => !p.FromSheet)
    applySort()
  }
  applyFilter()
  d3.select("#timeline-showold").on("change", applyFilter)

}

main()
