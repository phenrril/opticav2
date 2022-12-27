<?php
require_once "../conexion.php";
session_start();
$ojoD1= $_POST['ojoD1'];
$ojoD2= $_POST['ojoD2'];
$ojoD3= $_POST['ojoD3'];

$ojoI1= $_POST['ojoI1'];
$ojoI2= $_POST['ojoI2'];
$ojoI3= $_POST['ojoI3'];

$ojolD1= $_POST['ojoDl1'];
$ojolD2= $_POST['ojoDl2'];
$ojolD3= $_POST['ojoDl3'];

$ojolI1= $_POST['ojoIl1'];
$ojolI2= $_POST['ojoIl2'];
$ojolI3= $_POST['ojoIl3'];

$add= $_POST['add'];

$idgraduacion= 213;

$query = mysqli_query($conexion, ("INSERT INTO graduaciones (od_c_1, od_c_2, od_c_3, oi_c_1, oi_c_2, oi_c_3, od_l_1, od_l_2, od_l_3, oi_l_1, oi_l_2, oi_l_3, addg, id_paciente ) VALUES ($ojoD1 , $ojoD2, $ojoD3, $ojoI1, $ojoI2, $ojoI3, $ojolD1, $ojolD2, $ojolD3, $ojolI1, $ojolI2, $ojolI3, $add, $idgraduacion )"));


?>

