export type Point = { x: number; y: number };

function getAngle(center: Point, point: Point) {
  const x = point.x - center.x;
  const y = point.y - center.y;

  const angle = Math.atan2(y, x);

  if (angle <= 0) {
    return angle + Math.PI * 2;
  }

  return angle;
}

function getDistance(a: Point, b: Point) {
  return Math.sqrt((a.x - b.x) ** 2 + (a.y - b.y) ** 2);
}

function comparePoints(center: Point, a: Point, b: Point) {
  const angleA = getAngle(center, a);
  const angleB = getAngle(center, b);

  if (angleA < angleB) {
    return true;
  }

  const distanceA = getDistance(center, a);
  const distanceB = getDistance(center, b);

  if (angleA === angleB && distanceA < distanceB) {
    return true;
  }

  return false;
}

export function sortPoints(points: Point[], pointCenter: Point) {
  for (const point of points) {
    point.x -= pointCenter.x;
    point.y -= pointCenter.y;
  }

  points = points.sort((a, b) => {
    if (comparePoints({ x: 0, y: 0 }, a, b)) {
      return -1;
    }

    return 1;
  });

  for (const point of points) {
    point.x += pointCenter.x;
    point.y += pointCenter.y;
  }

  return points;
}
