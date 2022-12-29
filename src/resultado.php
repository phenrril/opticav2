<?php
require_once "../conexion.php";
session_start();

$ojolD1= $_POST['ojoDl1'];
if ($ojolD1 == ""){
    $ojolD1 = 0;}
$ojolD2= $_POST['ojoDl2'];
if ($ojolD2 == ""){
    $ojolD2 = 0;}
$ojolD3= $_POST['ojoDl3'];
if ($ojolD3 == ""){
    $ojolD3 = 0;}
$ojolI1= $_POST['ojoIl1'];
if ($ojolI1 == ""){
    $ojolI1 = 0;}
$ojolI2= $_POST['ojoIl2'];
if ($ojolI2 == ""){
    $ojolI2 = 0;}
$ojolI3= $_POST['ojoIl3'];
if ($ojolI3 == ""){
    $ojolI3 = 0;}
$ojoD1= $_POST['ojoD1'];
if ($ojoD1 == ""){
    $ojoD1 = 0;}
$ojoD2= $_POST['ojoD2'];
if ($ojoD2 == ""){
    $ojoD2 = 0;}
$ojoD3= $_POST['ojoD3'];
if ($ojoD3 == ""){
    $ojoD3 = 0;}
$ojoI1= $_POST['ojoI1'];
if ($ojoI1 == ""){
    $ojoI1 = 0;}
$ojoI2= $_POST['ojoI2'];
if ($ojoI2 == ""){
    $ojoI2 = 0;}
$ojoI3= $_POST['ojoI3'];
if ($ojoI3 == ""){
    $ojoI3 = 0;}

$add1= $_POST['add'];
if ($add1 == ""){
    $add1 = 0;}

$idventa1= mysqli_query($conexion, "SELECT MAX(id) as total FROM ventas");
$idventa2 = mysqli_fetch_assoc($idventa1);
$idventa3 = $idventa2['total'];
$idventa = $idventa3 + 1;

$obs = $_POST['obs'];
if ($obs == ""){
    $obs = "Sin Observaciones";
}

$query = mysqli_query($conexion, "INSERT INTO graduaciones(od_l_1, od_l_2, od_l_3, oi_l_1, oi_l_2, oi_l_3, od_c_1, od_c_2, od_c_3, oi_c_1, oi_c_2, oi_c_3, addg, id_venta, obs)  VALUES ('$ojolD1', '$ojolD2', '$ojolD3', '$ojolI1', '$ojolI2', '$ojolI3', '$ojoD1' , '$ojoD2', '$ojoD3', '$ojoI1', '$ojoI2', '$ojoI3', '$add1', '$idventa', '$obs')" );
if ($query) {    
    
    echo '<script>var ojoD1 = document.getElementById("ojoD1")</>';
    echo '<script>ojoD1.value = ""</script>';
    echo '<script>var ojoD2 = document.getElementById("ojoD2")</script>';
    echo '<script>ojoD2.value = ""</script>';
    echo '<script>var ojoD3 = document.getElementById("ojoD3")</script>';
    echo '<script>ojoD3.value = ""</script>';
    echo '<script>var ojoI1 = document.getElementById("ojoI1")</script>';
    echo '<script>ojoI1.value = ""</script>';
    echo '<script>var ojoI2 = document.getElementById("ojoI2")</script>';
    echo '<script>ojoI2.value = ""</script>';
    echo '<script>var ojoI3 = document.getElementById("ojoI3")</script>';
    echo '<script>ojoI3.value = ""</script>';
    echo '<script>var ojoDl1 = document.getElementById("ojoDl1")</script>';
    echo '<script>ojoDl1.value = ""</script>';
    echo '<script>var ojoDl2 = document.getElementById("ojoDl2")</script>';
    echo '<script>ojoDl2.value = ""</script>';
    echo '<script>var ojoDl3 = document.getElementById("ojoDl3")</script>';
    echo '<script>ojoDl3.value = ""</script>';
    echo '<script>var ojoIl1 = document.getElementById("ojoIl1")</script>';
    echo '<script>ojoIl1.value = ""</script>';
    echo '<script>var ojoIl2 = document.getElementById("ojoIl2")</script>';
    echo '<script>ojoIl2.value = ""</script>';
    echo '<script>var ojoIl3 = document.getElementById("ojoIl3")</script>';
    echo '<script>ojoIl3.value = ""</script>';
    echo '<script>var add1 = document.getElementById("add")</script>';
    echo '<script>add1.value = ""</script>';
    echo "<div class='alert alert-success'>Graduacion Agregada Correctamente</div>";
}
else {
    echo "<script>alert('Error al agregar Graduacion')</>";
}
?>

