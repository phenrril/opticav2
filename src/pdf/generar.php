<?php
require_once '../../conexion.php';
require_once 'fpdf/fpdf.php';
$pdf = new FPDF('P', 'mm', 'letter');
$pdf->AddPage();
$pdf->SetMargins(10, 10, 10);
$pdf->SetTitle("Ventas");
$pdf->SetFont('Arial', 'B', 12);
$id = $_GET['v'];
$idcliente = $_GET['cl'];
$config = mysqli_query($conexion, "SELECT * FROM configuracion");
//$consul = mysqli_query($conexion, "SELECT graduaciones.* FROM graduaciones JOIN ventas ON graduaciones.id_venta = ventas.id WHERE ventas.id = '$id'");
$gradu= mysqli_query($conexion, "SELECT * FROM graduaciones where id_venta='$id'");
$datos = mysqli_fetch_assoc($config);
$datos44 = mysqli_fetch_assoc($gradu);
$clientes = mysqli_query($conexion, "SELECT * FROM cliente WHERE idcliente = $idcliente");
$datosC = mysqli_fetch_assoc($clientes);
$ventas = mysqli_query($conexion, "SELECT d.*, p.codproducto, p.descripcion FROM detalle_venta d INNER JOIN producto p ON d.id_producto = p.codproducto WHERE d.id_venta = $id");
$pdf->Cell(195, 5, utf8_decode($datos['nombre']), 0, 1, 'C');
$pdf->Image("../../assets/img/logo.png", 180, 10, 30, 30, 'PNG');
$pdf->SetFont('Arial', 'B', 10);
$pdf->Cell(20, 5, utf8_decode("Teléfono: "), 0, 0, 'L');
$pdf->SetFont('Arial', '', 10);
$pdf->Cell(20, 5, $datos['telefono'], 0, 1, 'L');
$pdf->SetFont('Arial', 'B', 10);
$pdf->Cell(20, 5, utf8_decode("Dirección: "), 0, 0, 'L');
$pdf->SetFont('Arial', '', 10);
$pdf->Cell(20, 5, utf8_decode($datos['direccion']), 0, 1, 'L');
$pdf->SetFont('Arial', 'B', 10);
$pdf->Cell(20, 5, "Correo: ", 0, 0, 'L');
$pdf->SetFont('Arial', '', 10);
$pdf->Cell(20, 5, utf8_decode($datos['email']), 0, 1, 'L');
$pdf->Ln();
$pdf->SetFont('Arial', 'B', 10);
$pdf->SetFillColor(0, 0, 0);
$pdf->SetTextColor(255, 255, 255);
$pdf->Cell(196, 5, "Datos del cliente", 1, 1, 'C', 1);
$pdf->SetTextColor(0, 0, 0);
$pdf->Cell(90, 5, utf8_decode('Nombre'), 0, 0, 'L');
$pdf->Cell(50, 5, utf8_decode('Teléfono'), 0, 0, 'L');
$pdf->Cell(56, 5, utf8_decode('Dirección'), 0, 1, 'L');
$pdf->SetFont('Arial', '', 10);
$pdf->Cell(90, 5, utf8_decode($datosC['nombre']), 0, 0, 'L');
$pdf->Cell(50, 5, utf8_decode($datosC['telefono']), 0, 0, 'L');
$pdf->Cell(56, 5, utf8_decode($datosC['direccion']), 0, 1, 'L');
$pdf->Ln(3);
$pdf->SetFont('Arial', 'B', 10);
$pdf->SetTextColor(255, 255, 255);
$pdf->Cell(196, 5, "Detalle de Producto", 1, 1, 'C', 1);
$pdf->SetTextColor(0, 0, 0);
$pdf->Cell(14, 5, utf8_decode('N°'), 0, 0, 'L');
$pdf->Cell(90, 5, utf8_decode('Descripción'), 0, 0, 'L');
$pdf->Cell(25, 5, 'Cantidad', 0, 0, 'L');
$pdf->Cell(32, 5, 'Precio', 0, 0, 'L');
$pdf->Cell(35, 5, 'Sub Total.', 0, 1, 'L');
$pdf->SetFont('Arial', '', 10);
$contador = 1;
while ($row = mysqli_fetch_assoc($ventas)) {
    $pdf->Cell(14, 5, $contador, 0, 0, 'L');
    $pdf->Cell(90, 5, $row['descripcion'], 0, 0, 'L');
    $pdf->Cell(25, 5, $row['cantidad'], 0, 0, 'L');
    $pdf->Cell(32, 5, $row['precio'], 0, 0, 'L');
    $pdf->Cell(35, 5, number_format($row['cantidad'] * $row['precio'], 2, '.', ','), 0, 1, 'L');
    $contador++;
}
if ($datos44 == ""){

}else{

$pdf->SetFont('Arial', 'B', 10);
$pdf->SetTextColor(255, 255, 255);
$pdf->Cell(196, 5, "Graduaciones", 1, 1, 'C', 1);
$pdf->Ln(10);
$pdf->SetTextColor(0, 0, 0);

// Dibujar dos celdas en la misma fila para Ojo Derecho C y Ojo Derecho L
$pdf->Cell(45, 5, utf8_decode('Ojo Derecho C'), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode('Ojo Derecho L'), 1, 0, 'L');
$pdf->Ln(8);
$pdf->Cell(45, 5, utf8_decode($datos44['od_c_1']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['od_l_1']), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(45, 5, utf8_decode($datos44['od_c_2']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['od_l_2']), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(45, 5, utf8_decode($datos44['od_c_3']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['od_l_3']), 1, 0, 'L');
$pdf->Ln(8);

// Dibujar dos celdas en la misma fila para Ojo Izquierdo C y Ojo Izquierdo L
$pdf->Cell(45, 5, utf8_decode('Ojo Izquierdo C'), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode('Ojo Izquierdo L'), 1, 0, 'L');
$pdf->Ln(8);
$pdf->Cell(45, 5, utf8_decode($datos44['oi_c_1']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['oi_l_1']), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(45, 5, utf8_decode($datos44['oi_c_2']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['oi_l_2']), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(45, 5, utf8_decode($datos44['oi_c_3']), 1, 0, 'L');
$pdf->Cell(45, 5, utf8_decode($datos44['oi_l_3']), 1, 0, 'L');

$pdf->Ln(8);
$pdf->Cell(90, 5, utf8_decode('ADD'), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(90, 5, utf8_decode($datos44['addg']), 1, 0, 'L');
$pdf->Ln(8);
$pdf->Cell(90, 5, utf8_decode('Observaciones:'), 1, 0, 'L');
$pdf->Ln(6);
$pdf->Cell(90, 5, utf8_decode($datos44['obs']), 1, 0, 'L');
}
$pdf->Output("ventas.pdf", "I");

?>

