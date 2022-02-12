package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

const (
	// dimensions taken from elite c v2
	ecWidth  = 18.65
	ecLength = 33.0

	pinClearance  = 2.8
	pinOffset     = 0.15
	wallThickness = 1.0

	slotHeight  = 12.0
	slotWidth   = 28.0
	slotLength  = 3.9 + 2
	slotOffsetY = -2.0

	shieldHeight = 15.0
	shieldWidth  = 31.0
	shieldLength = 1.6

	trayHeight       = 5.0
	trayBottomHeight = 1.5
	ecTrayTranslateX = 3.70
)

func constructECTray() (sdf.SDF3, error) {
	ecTray, err := sdf.Box3D(sdf.V3{
		X: ecWidth + 2*wallThickness,
		Y: ecLength + 2*wallThickness,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	ecTray = sdf.Transform3D(ecTray, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: 0,
		Z: ecTray.BoundingBox().Size().Z / 2}))

	elitecBB, err := sdf.Box3D(sdf.V3{
		X: ecWidth,
		Y: ecLength,
		Z: trayHeight,
	}, 0)
	elitecBB = sdf.Transform3D(elitecBB, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: wallThickness,
		Z: elitecBB.BoundingBox().Size().Z/2 + trayBottomHeight}))

	if err != nil {
		return nil, err
	}
	ecTray = sdf.Difference3D(ecTray, elitecBB)

	ecTray = sdf.Transform3D(ecTray, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: -(ecLength + 2*wallThickness) / 2,
		Z: 0}))

	pinBox, err := sdf.Box3D(sdf.V3{
		X: pinClearance,
		Y: ecLength,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	pinBoxLeft := sdf.Transform3D(pinBox, sdf.Translate3d(sdf.V3{
		X: -ecWidth/2 + pinClearance/2,
		Y: wallThickness - pinBox.BoundingBox().Size().Y/2,
		Z: pinOffset + pinBox.BoundingBox().Size().Z/2}))

	ecTray = sdf.Difference3D(ecTray, pinBoxLeft)

	pinBoxRight := sdf.Transform3D(pinBox, sdf.Translate3d(sdf.V3{
		X: ecWidth/2 - pinClearance/2,
		Y: wallThickness - pinBox.BoundingBox().Size().Y/2,
		Z: pinOffset + pinBox.BoundingBox().Size().Z/2}))

	ecTray = sdf.Difference3D(ecTray, pinBoxRight)

	pinBoxBack, err := sdf.Box3D(sdf.V3{
		X: ecWidth - pinClearance/2,
		Y: pinClearance,
		Z: trayHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	pinBoxBack = sdf.Transform3D(pinBoxBack, sdf.Translate3d(sdf.V3{
		Y: -ecLength + wallThickness*2,
		Z: pinOffset + trayHeight/2}))

	ecTray = sdf.Difference3D(ecTray, pinBoxBack)

	pushHole, err := sdf.Cylinder3D(10, 2.5, 0)
	if err != nil {
		return nil, err
	}

	pushHole = sdf.Transform3D(pushHole, sdf.Translate3d(sdf.V3{
		Y: -ecLength/2 - ecLength/4,
		Z: 0}))

	ecTray = sdf.Difference3D(ecTray, pushHole)

	return ecTray, nil
}

func constructSlotAndShield() (sdf.SDF3, error) {
	slot, err := sdf.Box3D(sdf.V3{
		X: slotWidth,
		Y: slotLength,
		Z: slotHeight,
	}, 0)

	if err != nil {
		return nil, err
	}

	slot = sdf.Transform3D(slot, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: slotLength/2 + slotOffsetY,
		Z: slotHeight / 2}))

	tL := sdf.NewPolygon()
	tL.Add(0, 0)
	tL.Add(0, 2)
	tL.Add(-2, 0)
	tL.Close()

	tpL, err := sdf.Polygon2D(tL.Vertices())
	if err != nil {
		return nil, err
	}
	triangleL := sdf.Extrude3D(tpL, slotHeight)
	triangleL = sdf.Transform3D(triangleL, sdf.Translate3d(sdf.V3{
		X: -slotWidth / 2,
		Y: -triangleL.BoundingBox().Size().Y,
		Z: triangleL.BoundingBox().Size().Z / 2}))
	slot = sdf.Union3D(slot, triangleL)

	tR := sdf.NewPolygon()
	tR.Add(0, 0)
	tR.Add(0, -2)
	tR.Add(2, -2)
	tR.Close()

	tpR, err := sdf.Polygon2D(tR.Vertices())
	if err != nil {
		return nil, err
	}
	triangleR := sdf.Extrude3D(tpR, slotHeight)
	triangleR = sdf.Transform3D(triangleR, sdf.Translate3d(sdf.V3{
		X: +slotWidth / 2,
		Y: 0,
		Z: triangleR.BoundingBox().Size().Z / 2}))
	slot = sdf.Union3D(slot, triangleR)

	shield, err := sdf.Box3D(sdf.V3{
		X: shieldWidth,
		Y: shieldLength,
		Z: shieldHeight,
	}, 0)

	shield = sdf.Transform3D(shield, sdf.Translate3d(sdf.V3{
		X: 0,
		Y: shieldLength/2 + slotLength + slotOffsetY,
		Z: shieldHeight / 2}))

	slot = sdf.Union3D(slot, shield)

	return slot, nil
}

func holder() (sdf.SDF3, error) {
	//	trrsWidth := 6.0
	//	trrsLength := 13.30

	ecTray, err := constructECTray()

	// move it to the left side
	ecTray = sdf.Transform3D(ecTray, sdf.Translate3d(sdf.V3{
		X: ecTrayTranslateX,
		Y: 0,
		Z: 0}))

	slot, err := constructSlotAndShield()
	if err != nil {
		return nil, err
	}
	ecTray = sdf.Union3D(ecTray, slot)

	usbHole, err := createUSBCutout()
	if err != nil {
		return nil, err
	}

	ecTray = sdf.Difference3D(ecTray, usbHole)

	return ecTray, nil
}

func createUSBCutout() (sdf.SDF3, error) {
	slotCutoutThickness := 0.5

	usbCutout, err := sdf.Box3D(sdf.V3{
		X: 11.5,
		Y: slotLength + shieldLength - slotCutoutThickness,
		Z: 11 + trayBottomHeight,
	}, 2.5)

	if err != nil {
		return nil, err
	}

	usbCutout = sdf.Transform3D(usbCutout, sdf.Translate3d(sdf.V3{
		X: ecTrayTranslateX,
		Y: slotCutoutThickness + usbCutout.BoundingBox().Size().Y/2,
		Z: 0}))

	usbHole, err := sdf.Box3D(sdf.V3{
		X: 9.2,
		Y: slotLength + shieldLength + 6,
		Z: 3.3,
	}, 1.8)

	if err != nil {
		return nil, err
	}

	usbHole = sdf.Transform3D(usbHole, sdf.Translate3d(sdf.V3{
		X: ecTrayTranslateX,
		Y: usbHole.BoundingBox().Size().Y/2 - 3,
		Z: trayBottomHeight + usbHole.BoundingBox().Size().Z/2}))
	usbHole = sdf.Union3D(usbHole, usbCutout)

	return usbHole, nil
}

func main() {
	s, err := holder()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.RenderSTL(sdf.ScaleUniform3D(s, shrink), 300, "holder.stl")
}
