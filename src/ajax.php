<?php
require_once "../conexion.php";
session_start();

if (isset($_GET['q'])) {
    $datos = array();
    $nombre = $_GET['q'];
    $cliente = mysqli_query($conexion, "SELECT * FROM cliente WHERE nombre LIKE '%$nombre%' AND estado = 1");
    while ($row = mysqli_fetch_assoc($cliente)) {
        $data['id'] = $row['idcliente'];
        $data['label'] = $row['nombre'];
        $data['direccion'] = $row['direccion'];
        $data['telefono'] = $row['telefono'];
        $data['obrasocial'] = $row['obrasocial'];
        array_push($datos, $data);
    }
    echo json_encode($datos);
    die();
}else if (isset($_GET['pro'])) {
    $datos = array();
    $nombre = $_GET['pro'];
    //$producto = mysqli_query($conexion, "SELECT * FROM producto WHERE codigo LIKE '%" . $nombre . "%' OR descripcion LIKE '%" . $nombre . "%' AND estado = 1");
    $producto = mysqli_query($conexion, "SELECT * FROM producto WHERE codigo LIKE '%$nombre%'AND estado = 1 and existencia >0 OR descripcion LIKE '%$nombre%' AND estado =1 and existencia >0");
    while ($row = mysqli_fetch_assoc($producto)) {
        $data['id'] = $row['codproducto'];
        $data['label'] = $row['codigo'] . ' - ' .$row['descripcion'];
        $data['value'] = $row['descripcion'];
        $data['precio'] = $row['precio'];
        $data['existencia'] = $row['existencia'];
        array_push($datos, $data);
    }
    echo json_encode($datos);
    die();
}else if (isset($_GET['detalle'])) {
    $id = $_SESSION['idUser'];
    $datos = array();
    $detalle = mysqli_query($conexion, "SELECT d.*, p.codproducto, p.descripcion FROM detalle_temp d INNER JOIN producto p ON d.id_producto = p.codproducto WHERE d.id_usuario = $id");
    $sumar = mysqli_query($conexion, "SELECT total, SUM(total) AS total_pagar FROM detalle_temp WHERE id_usuario = $id");
    while ($row = mysqli_fetch_assoc($detalle)) {
        $data['id'] = $row['id'];
        $data['descripcion'] = $row['descripcion'];
        $data['cantidad'] = $row['cantidad'];
        $data['precio_venta'] = $row['precio_venta'];
        //$data['sub_total'] = number_format($row['precio_venta'] * $row['cantidad'], 2, '.', ',' ); // 
        $data['sub_total'] = $row['precio_venta'] * $row['cantidad'];
        array_push($datos, $data);
    }
    echo json_encode($datos);
    die();
} else if (isset($_GET['delete_detalle'])) {
    $id_detalle = $_GET['id'];
    $verificar = mysqli_query($conexion, "SELECT * FROM detalle_temp WHERE id = $id_detalle");
    $datos = mysqli_fetch_assoc($verificar);
    if ($datos['cantidad'] > 1) {
        $cantidad = $datos['cantidad'] - 1;
        $query = mysqli_query($conexion, "UPDATE detalle_temp SET cantidad = $cantidad WHERE id = $id_detalle");
        if ($query) {
            $msg = "restado";
        } else {
            $msg = "Error";
        }
    }else{
        $query = mysqli_query($conexion, "DELETE FROM detalle_temp WHERE id = $id_detalle");
        if ($query) {
            $msg = "ok";
        } else {
            $msg = "Error";
        }
    }
    echo $msg;
    die();
} else if (isset($_GET['procesarVenta'])) {
    $id_cliente = $_GET['id'];
    $id_user = $_SESSION['idUser'];
    $abona = $_GET['abona'];
    $fecha = date("Y-m-d");
    $restoant = $_GET['resto'];
    $descuento = $_GET['descuento'];    
    $obrasocial = $_GET['obrasocial'];
    if($obrasocial == ""){
        $obrasocial = 0;
    }
    $consulta = mysqli_query($conexion, "SELECT total, SUM(total) AS total_pagar FROM detalle_temp WHERE id_usuario = $id_user");
    $result = mysqli_fetch_assoc($consulta);   
    $total = ($result['total_pagar'] * $descuento) - $obrasocial; 
    $resto = $total - $abona - $obrasocial;
    $insertar = mysqli_query($conexion, "INSERT INTO ventas(id_cliente, total, id_usuario, abona, resto, obrasocial, fecha) VALUES ('$id_cliente', '$total', '$id_user', '$abona', '$resto', '$obrasocial', '$fecha')");
    if ($insertar) {
        $id_maximo = mysqli_query($conexion, "SELECT MAX(id) AS total FROM ventas");
        $resultId = mysqli_fetch_assoc($id_maximo);
        $ultimoId = $resultId['total'];
        $consultaDetalle = mysqli_query($conexion, "SELECT * FROM detalle_temp WHERE id_usuario = $id_user");
       $consultaDetalle2 = mysqli_query($conexion, "SELECT * FROM graduaciones_temp WHERE id_usuario = $id_user");
        while ($row = mysqli_fetch_assoc($consultaDetalle)) {
            $id_producto = $row['id_producto'];
            $cantidad = $row['cantidad'];
            $precio_original = $row['precio_venta'];
            $precio = $precio_original * $descuento;
            $insertarDet = mysqli_query($conexion, "INSERT INTO detalle_venta(id_producto, id_venta, cantidad, precio, precio_original, abona, resto, obrasocial ) VALUES ($id_producto, $ultimoId, $cantidad, '$precio', '$precio_original', '$abona', '$resto', '$obrasocial')");
            $postpagos = mysqli_query($conexion, "INSERT INTO postpagos(id_venta , abona, resto, precio, precio_original, id_cliente) VALUES ($ultimoId, '$abona', '$resto' , '$precio', '$precio_original', '$id_cliente')");
            $stockActual = mysqli_query($conexion, "SELECT * FROM producto WHERE codproducto = $id_producto");
            $stockNuevo = mysqli_fetch_assoc($stockActual);
            $stockTotal = $stockNuevo['existencia'] - $cantidad;
            $stock = mysqli_query($conexion, "UPDATE producto SET existencia = $stockTotal WHERE codproducto = $id_producto");
        } 
        while ($row2 = mysqli_fetch_assoc($consultaDetalle2)) {
            $ojolD1= $row2['od_l_1'];
            $ojolD2= $row2['od_l_2'];
            $ojolD3= $row2['od_l_3'];
            $ojolI1= $row2['oi_l_1'];
            $ojolI2= $row2['oi_l_2'];
            $ojolI3= $row2['oi_l_3'];
            $ojoD1= $row2['od_c_1'];
            $ojoD2= $row2['od_c_2'];
            $ojoD3= $row2['od_c_3'];
            $ojoI1= $row2['oi_c_1'];
            $ojoI2= $row2['oi_c_2'];
            $ojoI3= $row2['oi_c_3'];
            $add1= $row2['addg'];
            $obs = $row2['obs'];
            $id_venta2=$ultimoId;
            
        $insedsd = mysqli_query($conexion, "INSERT INTO graduaciones(od_l_1, od_l_2, od_l_3, oi_l_1, oi_l_2, oi_l_3, od_c_1, od_c_2, od_c_3, oi_c_1, oi_c_2, oi_c_3, addg, id_venta, obs)  VALUES ('$ojolD1', '$ojolD2', '$ojolD3', '$ojolI1', '$ojolI2', '$ojolI3', '$ojoD1' , '$ojoD2', '$ojoD3', '$ojoI1', '$ojoI2', '$ojoI3', '$add1', '$id_venta2', '$obs')" );
        }
        
        if ($insertarDet) {
            $eliminar = mysqli_query($conexion, "DELETE FROM detalle_temp WHERE id_usuario = $id_user");
            $eliminar = mysqli_query($conexion, "DELETE FROM descuento WHERE id_usuario = $id_user");
            $eliminar2 = mysqli_query($conexion, "DELETE FROM graduaciones_temp WHERE id_usuario = $id_user");
            $msg = array('id_cliente' => $id_cliente, 'id_venta' => $ultimoId);
        } 
    }else{
        $msg = array('mensaje' => 'error');
    }
    echo json_encode($msg);
    die();
}
if (isset($_POST['action'])) {
    $id = $_POST['id'];
    $cant = $_POST['cant'];
    $precio = $_POST['precio'];
    $id_user = $_SESSION['idUser'];
    $total = $precio * $cant;
    $verificar = mysqli_query($conexion, "SELECT * FROM detalle_temp WHERE id_producto = $id AND id_usuario = $id_user");
    $result = mysqli_num_rows($verificar);
    $datos = mysqli_fetch_assoc($verificar);
    if ($result > 0) {
        $cantidad = $datos['cantidad'] + 1;
        $total_precio = $cantidad * $total;
        $query = mysqli_query($conexion, "UPDATE detalle_temp SET cantidad = $cantidad, total = '$total_precio' WHERE id_producto = $id AND id_usuario = $id_user");
        if ($query) {
            $msg = "actualizado";
        } else {
            $msg = "Error al ingresar";
        }
    }else{
        $query = mysqli_query($conexion, "INSERT INTO detalle_temp(id_usuario, id_producto, cantidad, precio_venta, total) VALUES ($id_user, $id, $cant, '$precio', $total)");
        if ($query) {
            $msg = "registrado";
        }else{
            $msg = "Error al ingresar";
        }
    }
    echo json_encode($msg);
    die();
}
if (isset($_POST['cambio'])) {
    if (empty($_POST['actual']) || empty($_POST['nueva'])) {
        $msg = 'Los campos estan vacios';
    } else {
        $id = $_SESSION['idUser'];
        $actual = md5($_POST['actual']);
        $nueva = md5($_POST['nueva']);
        $consulta = mysqli_query($conexion, "SELECT * FROM usuario WHERE clave = '$actual' AND idusuario = $id");
        $result = mysqli_num_rows($consulta);
        if ($result == 1) {
            $query = mysqli_query($conexion, "UPDATE usuario SET clave = '$nueva' WHERE idusuario = $id");
            if ($query) {
                $msg = 'ok';
            }else{
                $msg = 'error';
            }
        } else {
            $msg = 'dif';
        }
        
    }
    echo $msg;
    die();
    
}

